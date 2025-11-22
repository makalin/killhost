# 🔪 KillHost
A tiny cross-platform terminal utility that **lists, tracks, kills, and opens localhost ports**.  
Perfect for cleaning up ghost processes left behind by tools like **Cursor**, **Vite**, **Next.js**, **Laravel**, **PHP local servers**, **Node dev servers**, and more.

KillHost gives you full control over your local ports without hunting for PIDs or using long system commands. Fast, simple, deadly.

---

## ✨ Features
- 🔍 **List all active localhost ports** (e.g., `:3000`, `:8000`, `:5173`, `:80`)  
- 🎯 **Kill processes by port number** (e.g., `killhost kill 3000`)  
- 👁️ **Track ports in real-time** (auto-refresh mode)  
- 🔗 **Open a running port in browser** (e.g., `killhost open 3000`)  
- ⚡ **Cross-platform** (macOS, Linux, Windows WSL)  
- 🧠 **Auto-detects orphaned dev servers** and offers to kill  
- 🪄 **Optional aliases** → `kh ls`, `kh kill 5173`, etc.

---

## 🚀 Installation
(Example for Go)

```sh
go install github.com/makalin/killhost@latest
````

Or clone:

```sh
git clone https://github.com/makalin/killhost
cd killhost
go build -o killhost
```

Move it into your PATH:

```sh
mv killhost /usr/local/bin/
```

---

## 🛠️ Usage

### ▶️ List all running localhost ports

```sh
killhost ls
```

Sample output:

```
:3000   Node (vite)     PID 4213
:5173   Vite Dev        PID 3892
:80     PHP httpd       PID 2783
```

---

### 💀 Kill a port

```sh
killhost kill 3000
```

Force kill:

```sh
killhost kill 3000 --force
```

---

### 🌐 Open a running port in browser

```sh
killhost open 5173
```

---

### ⏱️ Live watch mode (auto updates every 2s)

```sh
killhost watch
```

---

## ⚙️ Example Architecture

```
cmd/killhost
 ├─ main.go
 ├─ list.go
 ├─ kill.go
 ├─ watch.go
 └─ browser.go

internal/ports
 ├─ scanner.go   (lsof/netstat cross-platform logic)
 ├─ process.go
 └─ types.go
```

---

## 🧩 Roadmap

* Windows native support (without WSL)
* Auto-clean mode for stuck dev servers
* JSON output for scripting
* “Port rules” system (auto-kill on idle)
* GUI integration (menu bar indicator)

---

## 📜 License

MIT License

---

## 👤 Author

**Mehmet T. Akalın (makalin)**
[https://github.com/makalin](https://github.com/makalin)
