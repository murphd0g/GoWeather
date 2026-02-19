<script lang="ts">
  let address = $state("");
  let loading = $state(false);
  let result = $state("");
  let error = $state("");

  async function getWeather() {
    if (!address.trim()) {
      error = "Please enter an address";
      return;
    }

    loading = true;
    error = "";
    result = "";

    try {
      const response = await fetch(
        `/weather?address=${encodeURIComponent(address)}`,
      );
      const text = await response.text();

      if (response.ok) {
        // Check if it's a "No coordinates found" message
        if (text.startsWith("No coordinates found")) {
          error =
            text +
            "\n\nTip: Try being more specific (e.g., include city and state) or try a nearby landmark.";
        } else {
          result = text;
        }
      } else {
        error = text || "Failed to fetch weather data";
      }
    } catch (e) {
      error =
        "Network error. Make sure the Go backend is running on port 8080.";
    } finally {
      loading = false;
    }
  }
</script>

<svelte:head>
  <title>GoWeather</title>
</svelte:head>

<main>
  <div class="container">
    <h1>🌤️ GoWeather</h1>
    <p class="subtitle">Get weather forecasts for any US address</p>

    <form on:submit|preventDefault={getWeather}>
      <div class="input-group">
        <input
          type="text"
          bind:value={address}
          placeholder="Enter an address (e.g., 1600 Pennsylvania Ave NW, Washington, DC)"
          disabled={loading}
        />
        <button type="submit" disabled={loading}>
          {loading ? "🔄 Loading..." : "🔍 Get Weather"}
        </button>
      </div>
    </form>

    {#if error}
      <div class="error">
        <strong>❌ Error:</strong>
        {error}
      </div>
    {/if}

    {#if result}
      <div class="result">
        <h2>Weather Information</h2>
        <pre>{result}</pre>
      </div>
    {/if}
  </div>
</main>

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen,
      Ubuntu, Cantarell, sans-serif;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    min-height: 100vh;
    color: #333;
  }

  main {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    padding: 20px;
  }

  .container {
    background: white;
    border-radius: 20px;
    padding: 40px;
    max-width: 600px;
    width: 100%;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  }

  h1 {
    margin: 0 0 10px 0;
    font-size: 2.5rem;
    text-align: center;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .subtitle {
    text-align: center;
    color: #666;
    margin: 0 0 30px 0;
  }

  form {
    margin-bottom: 20px;
  }

  .input-group {
    display: flex;
    gap: 10px;
    flex-wrap: wrap;
  }

  input {
    flex: 1;
    min-width: 200px;
    padding: 14px 18px;
    border: 2px solid #e0e0e0;
    border-radius: 10px;
    font-size: 16px;
    transition: all 0.3s ease;
  }

  input:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  }

  input:disabled {
    background: #f5f5f5;
    cursor: not-allowed;
  }

  button {
    padding: 14px 28px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border: none;
    border-radius: 10px;
    font-size: 16px;
    font-weight: 600;
    cursor: pointer;
    transition:
      transform 0.2s ease,
      box-shadow 0.2s ease;
    white-space: nowrap;
  }

  button:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 8px 20px rgba(102, 126, 234, 0.4);
  }

  button:active:not(:disabled) {
    transform: translateY(0);
  }

  button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .error {
    background: #fee;
    border: 2px solid #fcc;
    border-radius: 10px;
    padding: 16px;
    color: #c33;
    margin-top: 20px;
  }

  .result {
    background: #f8f9fa;
    border: 2px solid #e9ecef;
    border-radius: 10px;
    padding: 20px;
    margin-top: 20px;
  }

  .result h2 {
    margin: 0 0 15px 0;
    color: #667eea;
    font-size: 1.5rem;
  }

  .result pre {
    margin: 0;
    white-space: pre-wrap;
    word-wrap: break-word;
    font-family: "Courier New", monospace;
    line-height: 1.6;
    color: #333;
  }

  @media (max-width: 600px) {
    .container {
      padding: 25px;
    }

    h1 {
      font-size: 2rem;
    }

    .input-group {
      flex-direction: column;
    }

    button {
      width: 100%;
    }
  }
</style>
