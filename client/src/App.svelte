<script lang="ts">
    let email = "";
    let password = "";
    let accessToken = "";
    let loading = false;
    let error = "";
    let output = "";

    const STORAGE_KEY = "memedroid_access_token";

    if (typeof localStorage !== "undefined") {
        accessToken = localStorage.getItem(STORAGE_KEY) ?? "";
    }

    async function login() {
        loading = true;
        error = "";
        output = "";

        try {
            const res = await fetch("/api/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ email, password }),
            });

            const text = await res.text();
            if (!res.ok) {
                error = text || `HTTP ${res.status}`;
                return;
            }

            output = text;
            const data = JSON.parse(text) as { access_token?: string };
            accessToken = data.access_token ?? "";
            localStorage.setItem(STORAGE_KEY, accessToken);
        } catch (e) {
            error = e instanceof Error ? e.message : "unknown error";
        } finally {
            loading = false;
        }
    }

    async function me() {
        loading = true;
        error = "";
        output = "";

        try {
            const res = await fetch("/api/me", {
                method: "GET",
                headers: {
                    Authorization: `Bearer ${accessToken}`,
                },
            });

            const text = await res.text();
            if (!res.ok) {
                error = text || `HTTP ${res.status}`;
                return;
            }

            output = text;
        } catch (e) {
            error = e instanceof Error ? e.message : "unknown error";
        } finally {
            loading = false;
        }
    }

    function clearToken() {
        accessToken = "";
        localStorage.removeItem(STORAGE_KEY);
    }
</script>

<main>
    <div class="container">
        <div class="banner">
            <strong>Under construction</strong>
            <span> — página de teste de login</span>
        </div>

        <h1>Cauldrun.Fun — Login Test</h1>
        <p class="hint">
            A API é acessada via <code>/api/*</code> (Caddy reverse proxy).
        </p>

        <section class="card">
            <h2>Login</h2>

            <form
                on:submit|preventDefault={() => {
                    void login();
                }}
            >
                <label>
                    <span>Email</span>
                    <input
                        type="email"
                        bind:value={email}
                        autocomplete="email"
                        required
                    />
                </label>

                <label>
                    <span>Senha</span>
                    <input
                        type="password"
                        bind:value={password}
                        autocomplete="current-password"
                        required
                    />
                </label>

                <div class="row">
                    <button type="submit" disabled={loading}>Entrar</button>
                    <button
                        type="button"
                        class="secondary"
                        on:click={() => void me()}
                        disabled={loading || !accessToken}
                    >
                        Testar /api/me
                    </button>
                </div>

                <div class="row">
                    <button
                        type="button"
                        class="danger"
                        on:click={clearToken}
                        disabled={loading || !accessToken}
                    >
                        Limpar token
                    </button>
                </div>
            </form>

            <div class="token">
                <div class="tokenHeader">
                    <span>Token (localStorage)</span>
                </div>
                <textarea readonly rows="3" value={accessToken}></textarea>
            </div>

            {#if error}
                <pre class="error">{error}</pre>
            {/if}

            {#if output}
                <pre class="output">{output}</pre>
            {/if}
        </section>
    </div>
</main>

<style>
    .container {
        max-width: 720px;
        margin: 0 auto;
        padding: 48px 16px;
    }

    .banner {
        display: inline-flex;
        gap: 8px;
        align-items: center;
        padding: 10px 12px;
        border-radius: 10px;
        background: rgba(255, 200, 0, 0.12);
        border: 1px solid rgba(255, 200, 0, 0.25);
        color: rgba(255, 255, 255, 0.92);
        margin-bottom: 18px;
    }

    h1 {
        margin: 0 0 8px;
        font-size: 2.2rem;
        letter-spacing: -0.02em;
    }

    .hint {
        margin: 0 0 20px;
        opacity: 0.72;
    }

    .card {
        padding: 18px;
        border-radius: 14px;
        background: rgba(255, 255, 255, 0.06);
        border: 1px solid rgba(255, 255, 255, 0.12);
    }

    .card h2 {
        margin: 0 0 12px;
        font-size: 1.1rem;
    }

    form {
        display: grid;
        gap: 12px;
    }

    label {
        display: grid;
        gap: 6px;
    }

    label span {
        font-size: 0.9rem;
        opacity: 0.8;
    }

    input {
        padding: 10px 12px;
        border-radius: 10px;
        border: 1px solid rgba(255, 255, 255, 0.16);
        background: rgba(0, 0, 0, 0.25);
        color: white;
        outline: none;
    }

    input:focus {
        border-color: rgba(120, 180, 255, 0.75);
    }

    .row {
        display: flex;
        gap: 10px;
        flex-wrap: wrap;
    }

    button {
        cursor: pointer;
        border-radius: 10px;
        border: 1px solid rgba(255, 255, 255, 0.16);
        padding: 10px 12px;
        color: white;
        background: rgba(120, 180, 255, 0.25);
    }

    button:disabled {
        cursor: not-allowed;
        opacity: 0.6;
    }

    .secondary {
        background: rgba(255, 255, 255, 0.08);
    }

    .danger {
        background: rgba(255, 80, 80, 0.2);
        border-color: rgba(255, 80, 80, 0.35);
    }

    .token {
        display: grid;
        gap: 8px;
        margin-top: 12px;
    }

    .tokenHeader {
        display: flex;
        justify-content: space-between;
        align-items: center;
        opacity: 0.8;
        font-size: 0.9rem;
    }

    textarea {
        width: 100%;
        padding: 10px 12px;
        border-radius: 10px;
        border: 1px solid rgba(255, 255, 255, 0.16);
        background: rgba(0, 0, 0, 0.25);
        color: white;
        resize: vertical;
    }

    pre {
        margin: 12px 0 0;
        padding: 12px;
        border-radius: 10px;
        overflow: auto;
        border: 1px solid rgba(255, 255, 255, 0.12);
        background: rgba(0, 0, 0, 0.25);
    }

    .error {
        border-color: rgba(255, 80, 80, 0.35);
    }
</style>
