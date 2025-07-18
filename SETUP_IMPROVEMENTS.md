# MediaManager Setup Improvements

This document outlines the improvements made to the MediaManager container setup to address common deployment issues.

## Changes Made

### 1. Automatic Directory Creation
- **Problem**: The application expected certain directories to exist but only created them in development mode
- **Solution**: The application now creates only the directories it manages (like `/app/images` for application data). Media directories are user-configured and mounted as volumes, so the application doesn't need to create them.

### 2. Images Directory Relocated
- **Problem**: Images were stored in `/data/images` which should be reserved for user media
- **Solution**: Images are now stored in `/app/images` within the application context, keeping user data separate from application data

### 3. Config Folder Approach
- **Problem**: The container required users to create a `config.toml` file before first boot, causing startup failures
- **Solution**: 
  - Changed from direct file mapping to config folder mapping
  - Added automatic config initialization on first boot
  - Example config is copied to the config folder if no config exists
  - Container can now start successfully without pre-existing configuration

## New Volume Structure

### Before:
```yaml
volumes:
  - ./data/:/data/
  - ./config.toml:/app/config.toml  # Required file that users had to create
```

### After:
```yaml
volumes:
  - ./data/:/data/
  - ./config/:/app/config/  # Config folder that gets auto-initialized
```

## First Boot Process

1. Container starts and checks for config directory
2. If config directory doesn't exist, it's created
3. If `config.toml` doesn't exist in the config directory, the example config is copied
4. Application-managed directories (like images) are created automatically
5. Database migrations run
6. **Default admin user is created if no users exist** (see Login Information below)
7. Application starts successfully

Note: Media directories are NOT created by the application - they should be mounted from your host system.

## Login Information

### Default Admin User
- If no users exist in the database, a default admin user is automatically created
- **Email**: First email from `admin_emails` in config.toml (default: `admin@example.com`)
- **Password**: `admin`
- **⚠️ IMPORTANT**: Change this password immediately after first login!

### Creating Additional Admin Users
- Register a new user with an email address listed in the `admin_emails` array in your config.toml
- Users with admin emails automatically become administrators upon registration

### Manual User Registration
- Access the web UI at `http://localhost:8000/`
- Look for registration/signup options in the interface
- Or use the API directly at `/api/v1/auth/register`

## User Experience Improvements

- **No pre-setup required**: Users can now run `docker-compose up` immediately
- **Clear configuration path**: Config file location is clearly communicated during startup
- **Example configuration**: Users get a working example config with helpful comments
- **Separation of concerns**: User data (`/data/`) is separate from app data (`/app/images`)

## Migration for Existing Users

If you have an existing setup:

1. Create a `config` folder in your docker-compose directory
2. Move your existing `config.toml` file into the `config` folder
3. Update your `docker-compose.yaml` to use the new volume mapping
4. Remove any existing `/data/images` folder (images will now be stored in the container)

## Directory Management

### Application-Managed Directories
- `/app/images`: Created automatically by the application for storing metadata images
- `/app/config`: Config folder mounted as volume, auto-initialized on first boot

### User-Managed Directories  
- Media directories (TV shows, movies, torrents): These are defined in your `config.toml` and should be mounted as volumes in your `docker-compose.yaml`
- The application does NOT create these directories - they should exist on your host system and be properly mounted

### Example Volume Mapping
```yaml
volumes:
  # Your actual media directories
  - /path/to/your/tv/shows:/data/tv
  - /path/to/your/movies:/data/movies
  # Config folder (auto-initialized)
  - ./config/:/app/config/
```

Then in your `config.toml`:
```toml
[[misc.tv_libraries]]
name = "TV Shows"
path = "/data/tv"  # This matches the container path from volume mount

[[misc.movie_libraries]]  
name = "Movies"
path = "/data/movies"  # This matches the container path from volume mount
```

## Environment Variables

- `CONFIG_DIR`: Path to config directory (default: `/app/config`)
- `MISC__IMAGE_DIRECTORY`: Path to images directory (default: `/app/images`)
- Media directory environment variables removed - these should be configured in `config.toml`