## 📄 Setting up Gemma 4 with Ollama

This guide covers how to install the backend engine and choose the right "brain" size for your Go Fiber API.

### 1. Install Ollama
If you haven't already, install the Ollama binary:
```bash
curl -fsSL https://ollama.com/install.sh | sh
```

### 2. Choose Your Gemma 4 Variant
Gemma 4 introduces **"Effective" (E)** parameters and **Mixture of Experts (MoE)** to run faster on consumer hardware like your Dell XPS.

| Model Tag | Size | Best For | Requirement |
| :--- | :--- | :--- | :--- |
| `gemma4:e2b` | 2.3B | **Extreme Speed.** Instant responses; great for background tasks. | 4GB RAM |
| `gemma4:e4b` | 4.5B | **The Sweet Spot.** Great for coding & reasoning on a laptop. | 8GB RAM |
| `gemma4:26b` | 26B (MoE) | **High Intelligence.** Uses MoE (only 4B active params). Best logic. | 16GB+ RAM |
| `gemma4:31b` | 31B | **Frontier Quality.** Best for complex creative writing/fine-tuning. | 24GB+ VRAM |

> **Recommendation for your XPS 15:** Use **`gemma4:e4b`** for general development or **`gemma4:26b`** if you want "Pro" level reasoning without killing your fans.

### 3. Pull the Model
Run the command for your chosen version:
```bash
# Example: Pulling the recommended laptop version
ollama pull gemma4:e4b
```

### 4. Update Your Go Code
If you decide to change the model (e.g. from `gemma4:e2b` to `gemma4:e4b`), update `internal/service/ollama.go`:

```go
chatPayload := models.OllamaChatRequest{
    Model: "gemma4:e4b", // Must match the 'ollama pull' tag
```

### 5. Verify the Local Link
Before running your Go app (`go run cmd/api/main.go`), ensure Ollama is visible:
```bash
curl http://localhost:11434/api/tags
```
If you see your model in that JSON list, you are ready to use the Web UI or API.

---

### Pro-Tip: Hardware Acceleration
Since you're on Linux (Ubuntu), ensure your NVIDIA drivers are active so Ollama can use your XPS's GPU (likely an RTX 3050/4050). Check this with:
```bash
nvidia-smi
```
If that command shows your GPU, Ollama will automatically offload Gemma 4 to the hardware, making your `/ask` endpoint feel nearly instantaneous.