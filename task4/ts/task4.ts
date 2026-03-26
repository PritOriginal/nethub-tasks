type Device = {
  id: number
  hostname: string
  ip: string
}

// Имитация запроса к API
async function fetchDevices(): Promise<Device[]> {
  const res = await fetch('/api/devices')

  if (!res.ok) {
    return Promise.reject(`HTTP error! status: ${res.status}`);
  }
  return (await res.json()) as Device[]
}

class DevicesStore {
  devices: Device[] = []
  isLoading = false

  loadAndFilterDevices = async (search: string): Promise<Device[]> => {
    this.isLoading = true
    try {
      const data = await fetchDevices()

      this.devices = data

      const normalizedSearch = search.trim().toLowerCase()
      const filtered = this.devices.filter((d) => d.hostname.toLowerCase().includes(normalizedSearch))

      return filtered

    } catch (error) {
      return []
    } finally {
      this.isLoading = false
    }
  }
}

// Пример использования (упрощённо)
async function example() {
  const deviceStore = new DevicesStore();

  const searchInput: HTMLInputElement | null = document.querySelector('#search')
  if (searchInput) {
    let debounceTimer: number | undefined;

    searchInput.oninput = () => {
      if (debounceTimer) {
        window.clearTimeout(debounceTimer);
      }

      debounceTimer = window.setTimeout(async () => {
        const list = await deviceStore.loadAndFilterDevices(searchInput.value);
        console.log('Devices:', list);
      }, 300);
    };
  }
}

example()

