# homeiota

Home IoT app

This repository is a monorepo for a home IoT monitoring and automation system. It includes backend services, device monitoring scripts, and a modern web frontend.

---

## Project Structure

- `go.alert.service/` — Go service for handling alerts and notifications. [See README](go.alert.service/README.md)
- `go.home.api/` — Go-based RESTful API for device data, including temperature, pump run times, and device heartbeats. [See README](go.home.api/README.md)
- `iot.pumpmon/` — Python-based monitor for pump devices. [See README](iot.pumpmon/README.md)
- `iot.tempmon/` — Python-based monitor for temperature sensors. [See README](iot.tempmon/README.md)
- `sveltekit.homeiota.app/` — SvelteKit web application for real-time device monitoring, visualization, and management. [See README](sveltekit.homeiota.app/README.md)

---

## Components

### go.alert.service
- Go service for sending alerts/notifications.
- Includes Dockerfile and Go module files.
- **Recommended run:**
  ```bash
  cd go.alert.service
  docker build -t homeiota-alert-service .
  docker run --env-file .env homeiota-alert-service
  ```
- [See go.alert.service README](go.alert.service/README.md)

### go.home.api
- RESTful API for device data and system integration.
- Includes Docker support and a `docker-compose.yml` for orchestration.
- **Recommended run:**
  ```bash
  cd go.home.api
  docker-compose up --build
  ```
- [See go.home.api README](go.home.api/README.md)

### iot.pumpmon
- Python script and supporting modules for monitoring pump devices.
- Configurable via `settings.toml`.
- [See iot.pumpmon README](iot.pumpmon/README.md)

### iot.tempmon
- Python script and supporting modules for monitoring temperature sensors.
- Configurable via `settings.toml`.
- [See iot.tempmon README](iot.tempmon/README.md)

### sveltekit.homeiota.app
- SvelteKit web frontend for device management and visualization.
- Uses Tailwind CSS, Prisma, and Vite.
- Contains static assets for PWA support.
- [See sveltekit.homeiota.app README](sveltekit.homeiota.app/README.md)


## Getting Started

Each component can be developed and deployed independently. See the respective subdirectory README files for setup and usage instructions.


## License

See [LICENSE](LICENSE) for details.
