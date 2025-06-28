import prettier from 'eslint-config-prettier';
import js from '@eslint/js';
import { includeIgnoreFile } from '@eslint/compat';
import svelte from 'eslint-plugin-svelte';
import globals from 'globals';
import { fileURLToPath } from 'node:url';
import ts from 'typescript-eslint';
import unusedImports from 'eslint-plugin-unused-imports';

const gitignorePath = fileURLToPath(new URL('./.gitignore', import.meta.url));

export default ts.config(
	{
		plugins: {
			'unused-imports': unusedImports
		},
		rules: {
			'@typescript-eslint/no-unused-vars': 'off',
			'unused-imports/no-unused-imports': 'error',
			'unused-imports/no-unused-vars': [
				'warn',
				{
					vars: 'all',
					varsIgnorePattern: '^_',
					args: 'after-used',
					argsIgnorePattern: '^_'
				}
			],
			'sort-imports': [
				'error',
				{
					ignoreDeclarationSort: true
				}
			]
		}
	},
	includeIgnoreFile(gitignorePath),
	js.configs.recommended,
	...ts.configs.recommended,
	...svelte.configs['flat/recommended'],
	prettier,
	...svelte.configs['flat/prettier'],
	{
		languageOptions: {
			globals: {
				...globals.browser,
				...globals.node
			}
		}
	},
	{
		files: ['**/*.svelte'],

		languageOptions: {
			parserOptions: {
				parser: ts.parser
			}
		}
	}
);
