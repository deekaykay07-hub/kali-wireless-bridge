# Colab Notebooks for Remote Unrestricted Models

This folder contains Google Colab notebooks for running powerful models remotely and exposing them to your Kali-Wireless TUI via ngrok.

## Main Notebook

- `kali_ollama_t4_unrestricted_15b_full.ipynb`
  - Uses the **Unrestricted 15B** model: `yqkm/Unrestricted-Knowledge-Will-Not-Refuse-15B-Q4_K_M-GGUF`
  - Designed for T4 GPU on Colab (free tier)
  - Sets up Ollama + pulls the GGUF + creates the exact model name the TUI expects
  - Starts ngrok tunnel on Ollama port
  - Prints the public URL to use as `OLLAMA_HOST`

## Quick Usage

1. Open the notebook in Google Colab.
2. **Runtime → Change runtime type → T4 GPU** (High-RAM if available).
3. Add your ngrok authtoken as a **Secret** named `NGROK_AUTH_TOKEN`.
4. Run all cells in order.
5. Copy the final public URL.
6. Use it in your TUI's `docker-compose.remote-colab.yml` (or set `OLLAMA_HOST` + `MISTRAL_MODEL`).

See the main TUI repo for `docker-compose.remote-colab.yml` example.

## Notes

- 15B models are heavy. Expect longer load times and possible OOM on low VRAM. Have fallback smaller models ready.
- Keep the Colab tab open or run the keep-alive cell.
- This setup lets you use much stronger reasoning than what fits on small VPS hardware while still using the Windows bridge for real hardware access.
