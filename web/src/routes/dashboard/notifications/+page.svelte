<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { env } from '$env/dynamic/public';
	import type { Notification } from '$lib/types';

	const apiUrl = env.PUBLIC_API_URL;

	interface NotificationResponse {
		id: string;
		read: boolean;
		message: string;
		timestamp: string;
	}

	let unreadNotifications: NotificationResponse[] = [];
	let readNotifications: NotificationResponse[] = [];
	let loading = true;
	let showRead = false;
	let markingAllAsRead = false;

	async function fetchNotifications() {
		try {
			loading = true;
			const [unreadResponse, allResponse] = await Promise.all([
				fetch(`${apiUrl}/notification/unread`, {
					method: 'GET',
					headers: {
						'Content-Type': 'application/json'
					},
					credentials: 'include'
				}),
				fetch(`${apiUrl}/notification`, {
					method: 'GET',
					headers: {
						'Content-Type': 'application/json'
					},
					credentials: 'include'
				})
			]);

			if (unreadResponse.ok) {
				unreadNotifications = await unreadResponse.json();
			}

			if (allResponse.ok) {
				const allNotifications: NotificationResponse[] = await allResponse.json();
				readNotifications = allNotifications.filter(n => n.read);
			}
		} catch (error) {
			console.error('Failed to fetch notifications:', error);
		} finally {
			loading = false;
		}
	}

	async function markAsRead(notificationId: string) {
		try {
			const response = await fetch(`${apiUrl}/notification/${notificationId}/read`, {
				method: 'PATCH',
				headers: {
					'Content-Type': 'application/json'
				},
				credentials: 'include'
			});

			if (response.ok) {
				// Move from unread to read
				const notification = unreadNotifications.find(n => n.id === notificationId);
				if (notification) {
					notification.read = true;
					readNotifications = [notification, ...readNotifications];
					unreadNotifications = unreadNotifications.filter(n => n.id !== notificationId);
				}
			}
		} catch (error) {
			console.error('Failed to mark notification as read:', error);
		}
	}

	async function markAsUnread(notificationId: string) {
		try {
			const response = await fetch(`${apiUrl}/notification/${notificationId}/unread`, {
				method: 'PATCH',
				headers: {
					'Content-Type': 'application/json'
				},
				credentials: 'include'
			});

			if (response.ok) {
				// Move from read to unread
				const notification = readNotifications.find(n => n.id === notificationId);
				if (notification) {
					notification.read = false;
					unreadNotifications = [notification, ...unreadNotifications];
					readNotifications = readNotifications.filter(n => n.id !== notificationId);
				}
			}
		} catch (error) {
			console.error('Failed to mark notification as unread:', error);
		}
	}

	async function deleteNotification(notificationId: string) {
		try {
			const response = await fetch(`${apiUrl}/notification/${notificationId}`, {
				method: 'DELETE',
				headers: {
					'Content-Type': 'application/json'
				},
				credentials: 'include'
			});

			if (response.ok) {
				// Remove from both lists
				unreadNotifications = unreadNotifications.filter(n => n.id !== notificationId);
				readNotifications = readNotifications.filter(n => n.id !== notificationId);
			}
		} catch (error) {
			console.error('Failed to delete notification:', error);
		}
	}

	async function markAllAsRead() {
		if (unreadNotifications.length === 0) return;

		try {
			markingAllAsRead = true;
			const promises = unreadNotifications.map(notification =>
				fetch(`${apiUrl}/notification/${notification.id}/read`, {
					method: 'PATCH',
					headers: {
						'Content-Type': 'application/json'
					},
					credentials: 'include'
				})
			);

			await Promise.all(promises);

			// Move all unread to read
			readNotifications = [...unreadNotifications.map(n => ({ ...n, read: true })), ...readNotifications];
			unreadNotifications = [];
		} catch (error) {
			console.error('Failed to mark all notifications as read:', error);
		} finally {
			markingAllAsRead = false;
		}
	}

	function formatTimestamp(timestamp: string): string {
		const date = new Date(timestamp);
		const now = new Date();
		const diffInMs = now.getTime() - date.getTime();
		const diffInMinutes = Math.floor(diffInMs / (1000 * 60));
		const diffInHours = Math.floor(diffInMs / (1000 * 60 * 60));
		const diffInDays = Math.floor(diffInMs / (1000 * 60 * 60 * 24));

		if (diffInMinutes < 1) return 'Just now';
		if (diffInMinutes < 60) return `${diffInMinutes}m ago`;
		if (diffInHours < 24) return `${diffInHours}h ago`;
		if (diffInDays < 7) return `${diffInDays}d ago`;

		return date.toLocaleDateString();
	}

	function getNotificationIcon(message: string): string {
		const lowerMessage = message.toLowerCase();

		if (lowerMessage.includes('downloaded') || lowerMessage.includes('successfully')) {
			return '✅';
		}
		if (lowerMessage.includes('error') || lowerMessage.includes('failed') || lowerMessage.includes('failure')) {
			return '❌';
		}
		if (lowerMessage.includes('missing') || lowerMessage.includes('not found')) {
			return '⚠️';
		}
		if (lowerMessage.includes('api')) {
			return '🔌';
		}
		if (lowerMessage.includes('indexer')) {
			return '🔍';
		}

		return '📢';
	}

	onMount(() => {
		fetchNotifications();

		// Refresh notifications every 30 seconds
		const interval = setInterval(fetchNotifications, 30000);
		return () => clearInterval(interval);
	});
</script>

<svelte:head>
	<title>Notifications - MediaManager</title>
</svelte:head>

<div class="container mx-auto px-4 py-8">
	<div class="flex justify-between items-center mb-6">
		<h1 class="text-3xl font-bold text-gray-900 dark:text-white">Notifications</h1>
		{#if unreadNotifications.length > 0}
			<button
				on:click={markAllAsRead}
				disabled={markingAllAsRead}
				class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
			>
				{#if markingAllAsRead}
					<div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
				{/if}
				Mark All as Read
			</button>
		{/if}
	</div>

	{#if loading}
		<div class="flex justify-center items-center py-12">
			<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
		</div>
	{:else}
		<!-- Unread Notifications -->
		<div class="mb-8">
			<div class="flex items-center gap-2 mb-4">
				<h2 class="text-xl font-semibold text-gray-900 dark:text-white">
					Unread Notifications
				</h2>
				{#if unreadNotifications.length > 0}
					<span class="bg-red-500 text-white text-xs px-2 py-1 rounded-full">
						{unreadNotifications.length}
					</span>
				{/if}
			</div>

			{#if unreadNotifications.length === 0}
				<div class="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg p-6 text-center">
					<div class="text-green-600 dark:text-green-400 text-4xl mb-2">✨</div>
					<p class="text-green-800 dark:text-green-200 font-medium">All caught up!</p>
					<p class="text-green-600 dark:text-green-400 text-sm">No unread notifications</p>
				</div>
			{:else}
				<div class="space-y-3">
					{#each unreadNotifications as notification (notification.id)}
						<div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4 shadow-sm">
							<div class="flex items-start justify-between gap-4">
								<div class="flex items-start gap-3 flex-1">
									<span class="text-2xl">{getNotificationIcon(notification.message)}</span>
									<div class="flex-1">
										<p class="text-gray-900 dark:text-white font-medium">
											{notification.message}
										</p>
										<p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
											{formatTimestamp(notification.timestamp)}
										</p>
									</div>
								</div>
								<div class="flex items-center gap-2">
									<button
										on:click={() => markAsRead(notification.id)}
										class="p-2 text-blue-600 hover:bg-blue-100 dark:hover:bg-blue-800 rounded-lg transition-colors"
										title="Mark as read"
									>
										<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
										</svg>
									</button>
									<button
										on:click={() => deleteNotification(notification.id)}
										class="p-2 text-red-600 hover:bg-red-100 dark:hover:bg-red-800 rounded-lg transition-colors"
										title="Delete notification"
									>
										<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
										</svg>
									</button>
								</div>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>

		<!-- Read Notifications Toggle -->
		<div class="mb-4">
			<button
				on:click={() => showRead = !showRead}
				class="flex items-center gap-2 text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white transition-colors"
			>
				<svg
					class="w-4 h-4 transition-transform {showRead ? 'rotate-90' : ''}"
					fill="none"
					stroke="currentColor"
					viewBox="0 0 24 24"
				>
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
				</svg>
				<span>Read Notifications ({readNotifications.length})</span>
			</button>
		</div>

		<!-- Read Notifications -->
		{#if showRead}
			<div>
				{#if readNotifications.length === 0}
					<div class="bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 text-center">
						<p class="text-gray-500 dark:text-gray-400">No read notifications</p>
					</div>
				{:else}
					<div class="space-y-3">
						{#each readNotifications as notification (notification.id)}
							<div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-4 shadow-sm opacity-75">
								<div class="flex items-start justify-between gap-4">
									<div class="flex items-start gap-3 flex-1">
										<span class="text-2xl opacity-50">{getNotificationIcon(notification.message)}</span>
										<div class="flex-1">
											<p class="text-gray-700 dark:text-gray-300">
												{notification.message}
											</p>
											<p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
												{formatTimestamp(notification.timestamp)}
											</p>
										</div>
									</div>
									<div class="flex items-center gap-2">
										<button
											on:click={() => markAsUnread(notification.id)}
											class="p-2 text-blue-600 hover:bg-blue-100 dark:hover:bg-blue-800 rounded-lg transition-colors"
											title="Mark as unread"
										>
											<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 7.89a2 2 0 002.83 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
											</svg>
										</button>
										<button
											on:click={() => deleteNotification(notification.id)}
											class="p-2 text-red-600 hover:bg-red-100 dark:hover:bg-red-800 rounded-lg transition-colors"
											title="Delete notification"
										>
											<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
											</svg>
										</button>
									</div>
								</div>
							</div>
						{/each}
					</div>
				{/if}
			</div>
		{/if}
	{/if}
</div>
