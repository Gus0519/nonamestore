<script>
  import { onMount } from 'svelte';

  const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

  let products = [];
  let cart = [];
  let loading = true;
  let error = null;

  onMount(async () => {
    try {
      const response = await fetch(`${API_URL}/products`);
      if (!response.ok) throw new Error('No se pudo conectar con el Backend en Go');
      products = await response.json();
      loading = false;
    } catch (err) {
      error = err.message;
      loading = false;
    }
  });

  function addToCart(product) {
    const existingIndex = cart.findIndex(item => item.id === product.id);
    if (existingIndex !== -1) {
      cart[existingIndex].quantity += 1;
      cart = [...cart];
    } else {
      cart = [...cart, { ...product, quantity: 1 }];
    }
  }

  function removeFromCart(productId) {
    cart = cart.filter(item => item.id !== productId);
  }

  $: cartTotal = cart.reduce((sum, item) => sum + (item.price * item.quantity), 0);
  $: cartCount = cart.reduce((sum, item) => sum + item.quantity, 0);
</script>

<main class="app-container">
  <header class="navbar">
    <div class="brand">🚀 NoNameStore <span class="badge">MVP</span></div>
    <div class="cart-status">🛒 Carrito ({cartCount})</div>
  </header>

  <div class="layout">
    <section class="catalog-section">
      <h2>Catálogo de Productos</h2>
      
      {#if loading}
        <p class="status-msg">⏳ Cargando vitrina desde el motor de Go...</p>
      {:else if error}
        <p class="status-msg error">❌ Error: {error}. ¿Levantaste el backend?</p>
      {:else}
        <div class="product-grid">
          {#each products as product}
            <div class="product-card">
              <img src={product.image_url} alt={product.name} class="product-img" />
              <div class="product-info">
                <h3>{product.name}</h3>
                <p class="desc">{product.description}</p>
                <div class="card-footer">
                  <span class="price">${product.price.toFixed(2)}</span>
                  <button on:click={() => addToCart(product)}>Agregar</button>
                </div>
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </section>

    <aside class="cart-sidebar">
      <h2>Tu Pedido</h2>
      {#if cart.length === 0}
        <p class="empty-cart">El carrito está vacío. Elige hardware de verdad.</p>
      {:else}
        <div class="cart-list">
          {#each cart as item}
            <div class="cart-item">
              <div class="item-details">
                <h4>{item.name} (x{item.quantity})</h4>
                <p class="item-subtotal">${(item.price * item.quantity).toFixed(2)}</p>
              </div>
              <button class="btn-remove" on:click={() => removeFromCart(item.id)}>❌</button>
            </div>
          {/each}
        </div>
        <div class="cart-summary">
          <h3>Total: ${cartTotal.toFixed(2)}</h3>
          <button class="btn-checkout" on:click={() => alert('¡MVP Completo! Próximo paso: Kubernetes')}>Proceder al Pago</button>
        </div>
      {/if}
    </aside>
  </div>
</main>

<style>
  :global(body) {
    margin: 0;
    font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    background-color: #121214;
    color: #e1e1e6;
  }

  .app-container {
    max-width: 1400px;
    margin: 0 auto;
    padding: 20px;
  }

  .navbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 15px 25px;
    background: #202024;
    border-radius: 8px;
    margin-bottom: 30px;
    border-bottom: 2px solid #00b37e;
  }

  .brand {
    font-size: 1.4rem;
    font-weight: bold;
  }

  .badge {
    font-size: 0.75rem;
    background: #00b37e;
    color: #121214;
    padding: 2px 6px;
    border-radius: 4px;
    vertical-align: middle;
  }

  .cart-status {
    font-size: 1.1rem;
    font-weight: 500;
  }

  .layout {
    display: grid;
    grid-template-columns: 3fr 1fr;
    gap: 30px;
  }

  @media (max-width: 900px) {
    .layout {
      grid-template-columns: 1fr;
    }
  }

  h2 {
    margin-top: 0;
    border-left: 4px solid #00b37e;
    padding-left: 10px;
  }

  .status-msg {
    text-align: center;
    padding: 40px;
    background: #202024;
    border-radius: 8px;
  }

  .status-msg.error {
    color: #f75a68;
    border: 1px solid #f75a68;
  }

  .product-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 20px;
  }

  .product-card {
    background: #202024;
    border-radius: 8px;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    transition: transform 0.2s;
  }

  .product-card:hover {
    transform: translateY(-5px);
  }

  .product-img {
    width: 100%;
    height: 180px;
    object-fit: cover;
  }

  .product-info {
    padding: 15px;
    display: flex;
    flex-direction: column;
    flex-grow: 1;
  }

  .product-info h3 {
    margin: 0 0 10px 0;
    font-size: 1.2rem;
  }

  .desc {
    font-size: 0.9rem;
    color: #8d8d99;
    flex-grow: 1;
    margin-bottom: 15px;
  }

  .card-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .price {
    font-size: 1.3rem;
    font-weight: bold;
    color: #00b37e;
  }

  button {
    background: #00b37e;
    color: #121214;
    border: none;
    padding: 8px 16px;
    border-radius: 4px;
    font-weight: bold;
    cursor: pointer;
    transition: background 0.2s;
  }

  button:hover {
    background: #008f64;
  }

  .cart-sidebar {
    background: #202024;
    padding: 20px;
    border-radius: 8px;
    height: fit-content;
  }

  .empty-cart {
    color: #8d8d99;
    font-style: italic;
    font-size: 0.95rem;
  }

  .cart-list {
    margin-bottom: 20px;
  }

  .cart-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 0;
    border-bottom: 1px solid #323238;
  }

  .item-details h4 {
    margin: 0;
    font-size: 0.95rem;
  }

  .item-subtotal {
    margin: 5px 0 0 0;
    font-size: 0.9rem;
    color: #00b37e;
  }

  .btn-remove {
    background: transparent;
    padding: 4px;
    font-size: 0.8rem;
  }

  .btn-remove:hover {
    background: #323238;
  }

  .cart-summary {
    border-top: 2px solid #323238;
    padding-top: 15px;
  }

  .btn-checkout {
    width: 100%;
    padding: 12px;
    font-size: 1rem;
    margin-top: 10px;
  }
</style>