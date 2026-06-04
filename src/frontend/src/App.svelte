<script lang="ts">
  import { onMount } from 'svelte';

  const API_URL = import.meta.env.VITE_API_URL || '';

  interface Product {
    id: number;
    name: string;
    scientific_name?: string;
    description: string;
    price: number;
    image_url: string;
    care_level?: string;
  }

  interface CartItem extends Product {
    quantity: number;
  }

  let products: Product[] = $state([]);
  let cart: CartItem[] = $state([]);
  let loading = $state(true);
  let error: string | null = $state(null);
  let cartOpen = $state(false);

  // Productos de demostración para cuando no hay backend
  const demoProducts: Product[] = [
    {
      id: 1,
      name: 'Monstera Deliciosa',
      scientific_name: 'Monstera deliciosa',
      description: 'La iconica planta de hojas perforadas que transforma cualquier espacio en un oasis tropical.',
      price: 89.00,
      image_url: '/monstera.png',
      care_level: 'Facil'
    },
    {
      id: 2,
      name: 'Ficus Lyrata',
      scientific_name: 'Ficus lyrata',
      description: 'Elegante arbol de interior con hojas grandes en forma de violin, perfecto para espacios luminosos.',
      price: 120.00,
      image_url: '/fiddle-leaf.png',
      care_level: 'Moderado'
    },
    {
      id: 3,
      name: 'Sansevieria',
      scientific_name: 'Dracaena trifasciata',
      description: 'Resistente y purificadora de aire, ideal para principiantes. Requiere muy poco mantenimiento.',
      price: 45.00,
      image_url: '/snake-plant.png',
      care_level: 'Muy Facil'
    }
  ];

  onMount(async () => {
    try {
      const response = await fetch(`${API_URL}/products`);
      if (!response.ok) throw new Error('No se pudo conectar con el Backend en Go');
      products = await response.json();
      loading = false;
    } catch (err) {
      // Usar productos demo si el backend no esta disponible
      products = demoProducts;
      error = null;
      loading = false;
    }
  });

  function addToCart(product: Product) {
    const existingIndex = cart.findIndex(item => item.id === product.id);
    if (existingIndex !== -1) {
      cart[existingIndex].quantity += 1;
      cart = [...cart];
    } else {
      cart = [...cart, { ...product, quantity: 1 }];
    }
  }

  function removeFromCart(productId: number) {
    cart = cart.filter(item => item.id !== productId);
  }

  function updateQuantity(productId: number, delta: number) {
    const index = cart.findIndex(item => item.id === productId);
    if (index !== -1) {
      cart[index].quantity += delta;
      if (cart[index].quantity <= 0) {
        cart = cart.filter(item => item.id !== productId);
      } else {
        cart = [...cart];
      }
    }
  }

  let cartTotal = $derived(cart.reduce((sum, item) => sum + (item.price * item.quantity), 0));
  let cartCount = $derived(cart.reduce((sum, item) => sum + item.quantity, 0));
</script>

<main class="app-container">
  <!-- Navbar -->
  <header class="navbar">
    <div class="nav-content">
      <a href="/" class="brand">
        <svg class="logo-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M12 2C8 6 4 10 4 14c0 4.4 3.6 8 8 8s8-3.6 8-8c0-4-4-8-8-12z" />
          <path d="M12 22V8" />
          <path d="M8 14c2-2 4-2 4-2s2 0 4 2" />
        </svg>
        <span class="brand-text">NoNameStore</span>
      </a>

      <nav class="nav-links">
        <a href="#catalogo">Catalogo</a>
        <a href="#cuidados">Cuidados</a>
        <a href="#nosotros">Nosotros</a>
      </nav>

      <div class="nav-actions">
        <button class="icon-btn" aria-label="Buscar">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8" />
            <path d="M21 21l-4.35-4.35" />
          </svg>
        </button>
        <button class="icon-btn" aria-label="Cuenta">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
            <circle cx="12" cy="7" r="4" />
          </svg>
        </button>
        <button class="cart-btn" onclick={() => cartOpen = true} aria-label="Carrito">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z" />
            <line x1="3" y1="6" x2="21" y2="6" />
            <path d="M16 10a4 4 0 0 1-8 0" />
          </svg>
          {#if cartCount > 0}
            <span class="cart-badge">{cartCount}</span>
          {/if}
        </button>
      </div>
    </div>
  </header>

  <!-- Hero Section -->
  <section class="hero">
    <div class="hero-content">
      <p class="hero-tagline">Naturaleza en tu hogar</p>
      <h1 class="hero-title">Transforma tu espacio con plantas de interior</h1>
      <p class="hero-description">
        Descubre nuestra coleccion curada de plantas que purifican el aire
        y llenan de vida cada rincon de tu hogar.
      </p>
      <a href="#catalogo" class="btn-primary">Explorar Catalogo</a>
    </div>
  </section>

  <!-- Catalog Section -->
  <section class="catalog-section" id="catalogo">
    <div class="section-header">
      <h2>Coleccion Destacada</h2>
      <p>Plantas seleccionadas para cada tipo de espacio y nivel de experiencia</p>
    </div>

    {#if loading}
      <div class="status-container">
        <div class="loading-spinner"></div>
        <p class="status-msg">Cargando catalogo desde el backend Go...</p>
      </div>
    {:else if error}
      <div class="status-container error">
        <p class="status-msg">{error}</p>
        <p class="status-hint">Asegurate de que el backend este corriendo</p>
      </div>
    {:else}
      <div class="product-grid">
        {#each products as product}
          <article class="product-card">
            <div class="product-image-container">
              <img src={product.image_url} alt={product.name} class="product-img" />
              {#if product.care_level}
                <span class="care-badge">{product.care_level}</span>
              {/if}
            </div>
            <div class="product-info">
              <div class="product-header">
                <h3>{product.name}</h3>
                {#if product.scientific_name}
                  <p class="scientific-name">{product.scientific_name}</p>
                {/if}
              </div>
              <p class="product-desc">{product.description}</p>
              <div class="card-footer">
                <span class="price">${product.price.toFixed(2)}</span>
                <button class="btn-add" onclick={() => addToCart(product)}>
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="12" y1="5" x2="12" y2="19" />
                    <line x1="5" y1="12" x2="19" y2="12" />
                  </svg>
                  Agregar
                </button>
              </div>
            </div>
          </article>
        {/each}
      </div>
    {/if}
  </section>

  <!-- Features Section -->
  <section class="features-section">
    <div class="feature">
      <div class="feature-icon">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <rect x="1" y="3" width="15" height="13" rx="2" />
          <path d="M16 8h4a2 2 0 0 1 2 2v7a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2v-2" />
        </svg>
      </div>
      <h3>Envio Protegido</h3>
      <p>Empaque especial para que tu planta llegue perfecta</p>
    </div>
    <div class="feature">
      <div class="feature-icon">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
          <path d="M9 12l2 2 4-4" />
        </svg>
      </div>
      <h3>Garantia de Salud</h3>
      <p>30 dias de garantia en todas nuestras plantas</p>
    </div>
    <div class="feature">
      <div class="feature-icon">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z" />
        </svg>
      </div>
      <h3>Soporte Experto</h3>
      <p>Asesoria personalizada para el cuidado de tus plantas</p>
    </div>
  </section>
</main>

<!-- Cart Sidebar Overlay -->
{#if cartOpen}
  <div class="cart-overlay" onclick={() => cartOpen = false} role="button" tabindex="-1" aria-label="Cerrar carrito"></div>
{/if}

<!-- Cart Sidebar -->
<aside class="cart-sidebar" class:open={cartOpen}>
  <div class="cart-header">
    <h2>Tu Carrito</h2>
    <button class="close-btn" onclick={() => cartOpen = false} aria-label="Cerrar">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <line x1="18" y1="6" x2="6" y2="18" />
        <line x1="6" y1="6" x2="18" y2="18" />
      </svg>
    </button>
  </div>

  {#if cart.length === 0}
    <div class="empty-cart">
      <svg class="empty-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z" />
        <line x1="3" y1="6" x2="21" y2="6" />
        <path d="M16 10a4 4 0 0 1-8 0" />
      </svg>
      <p>Tu carrito esta vacio</p>
      <span>Agrega plantas para comenzar</span>
    </div>
  {:else}
    <div class="cart-items">
      {#each cart as item}
        <div class="cart-item">
          <img src={item.image_url} alt={item.name} class="cart-item-img" />
          <div class="cart-item-details">
            <h4>{item.name}</h4>
            <p class="cart-item-price">${item.price.toFixed(2)}</p>
            <div class="quantity-controls">
              <button onclick={() => updateQuantity(item.id, -1)} aria-label="Reducir cantidad">-</button>
              <span>{item.quantity}</span>
              <button onclick={() => updateQuantity(item.id, 1)} aria-label="Aumentar cantidad">+</button>
            </div>
          </div>
          <button class="remove-btn" onclick={() => removeFromCart(item.id)} aria-label="Eliminar">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M3 6h18" />
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
            </svg>
          </button>
        </div>
      {/each}
    </div>

    <div class="cart-footer">
      <div class="cart-total">
        <span>Subtotal</span>
        <span class="total-amount">${cartTotal.toFixed(2)}</span>
      </div>
      <button class="btn-checkout">Proceder al Pago</button>
      <p class="shipping-note">Envio calculado en el siguiente paso</p>
    </div>
  {/if}
</aside>

<style>
  :global(*) {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  :global(body) {
    font-family: 'DM Sans', system-ui, -apple-system, sans-serif;
    background-color: #0d1210;
    color: #e8ede9;
    line-height: 1.6;
  }

  .app-container {
    min-height: 100vh;
  }

  /* Navbar */
  .navbar {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    z-index: 100;
    background: rgba(13, 18, 16, 0.95);
    backdrop-filter: blur(12px);
    border-bottom: 1px solid rgba(58, 90, 64, 0.3);
  }

  .nav-content {
    max-width: 1280px;
    margin: 0 auto;
    padding: 1rem 2rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .brand {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    text-decoration: none;
    color: inherit;
  }

  .logo-icon {
    width: 32px;
    height: 32px;
    color: #5a8a62;
  }

  .brand-text {
    font-family: 'Playfair Display', serif;
    font-size: 1.5rem;
    font-weight: 600;
    letter-spacing: -0.02em;
  }

  .nav-links {
    display: flex;
    gap: 2.5rem;
  }

  .nav-links a {
    color: #9ca89e;
    text-decoration: none;
    font-size: 0.9rem;
    font-weight: 500;
    transition: color 0.2s;
  }

  .nav-links a:hover {
    color: #e8ede9;
  }

  .nav-actions {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .icon-btn, .cart-btn {
    background: transparent;
    border: none;
    color: #9ca89e;
    padding: 0.5rem;
    cursor: pointer;
    border-radius: 8px;
    transition: all 0.2s;
    position: relative;
  }

  .icon-btn:hover, .cart-btn:hover {
    background: rgba(90, 138, 98, 0.15);
    color: #e8ede9;
  }

  .icon-btn svg, .cart-btn svg {
    width: 22px;
    height: 22px;
  }

  .cart-badge {
    position: absolute;
    top: 0;
    right: 0;
    background: #5a8a62;
    color: #0d1210;
    font-size: 0.7rem;
    font-weight: 700;
    width: 18px;
    height: 18px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  /* Hero */
  .hero {
    padding: 10rem 2rem 6rem;
    max-width: 1280px;
    margin: 0 auto;
  }

  .hero-content {
    max-width: 640px;
  }

  .hero-tagline {
    color: #5a8a62;
    font-size: 0.875rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    margin-bottom: 1rem;
  }

  .hero-title {
    font-family: 'Playfair Display', serif;
    font-size: 3.5rem;
    font-weight: 600;
    line-height: 1.1;
    margin-bottom: 1.5rem;
    letter-spacing: -0.02em;
  }

  .hero-description {
    color: #9ca89e;
    font-size: 1.125rem;
    margin-bottom: 2rem;
    max-width: 480px;
  }

  .btn-primary {
    display: inline-flex;
    align-items: center;
    padding: 1rem 2rem;
    background: #5a8a62;
    color: #0d1210;
    text-decoration: none;
    font-weight: 600;
    border-radius: 8px;
    transition: all 0.2s;
  }

  .btn-primary:hover {
    background: #4a7a52;
    transform: translateY(-2px);
  }

  /* Catalog Section */
  .catalog-section {
    padding: 4rem 2rem;
    max-width: 1280px;
    margin: 0 auto;
  }

  .section-header {
    margin-bottom: 3rem;
  }

  .section-header h2 {
    font-family: 'Playfair Display', serif;
    font-size: 2.25rem;
    font-weight: 600;
    margin-bottom: 0.5rem;
  }

  .section-header p {
    color: #9ca89e;
    font-size: 1rem;
  }

  .status-container {
    text-align: center;
    padding: 4rem 2rem;
    background: rgba(26, 35, 30, 0.5);
    border-radius: 12px;
    border: 1px solid rgba(58, 90, 64, 0.2);
  }

  .loading-spinner {
    width: 40px;
    height: 40px;
    border: 3px solid rgba(90, 138, 98, 0.2);
    border-top-color: #5a8a62;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin: 0 auto 1rem;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .status-msg {
    color: #9ca89e;
  }

  .status-container.error {
    border-color: rgba(247, 90, 104, 0.3);
  }

  .status-container.error .status-msg {
    color: #f75a68;
  }

  .status-hint {
    color: #6b736c;
    font-size: 0.875rem;
    margin-top: 0.5rem;
  }

  /* Product Grid */
  .product-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 2rem;
  }

  @media (max-width: 1024px) {
    .product-grid {
      grid-template-columns: repeat(2, 1fr);
    }
  }

  @media (max-width: 640px) {
    .product-grid {
      grid-template-columns: 1fr;
    }
  }

  /* Product Card */
  .product-card {
    background: #1a231e;
    border-radius: 16px;
    overflow: hidden;
    border: 1px solid rgba(58, 90, 64, 0.2);
    transition: all 0.3s ease;
  }

  .product-card:hover {
    transform: translateY(-8px);
    border-color: rgba(90, 138, 98, 0.4);
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
  }

  .product-image-container {
    position: relative;
    aspect-ratio: 4/3;
    overflow: hidden;
    background: #141a17;
  }

  .product-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.4s ease;
  }

  .product-card:hover .product-img {
    transform: scale(1.05);
  }

  .care-badge {
    position: absolute;
    top: 1rem;
    left: 1rem;
    background: rgba(13, 18, 16, 0.85);
    color: #5a8a62;
    padding: 0.375rem 0.75rem;
    border-radius: 20px;
    font-size: 0.75rem;
    font-weight: 600;
    backdrop-filter: blur(8px);
  }

  .product-info {
    padding: 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .product-header h3 {
    font-size: 1.25rem;
    font-weight: 600;
    margin-bottom: 0.25rem;
  }

  .scientific-name {
    font-style: italic;
    color: #6b736c;
    font-size: 0.875rem;
  }

  .product-desc {
    color: #9ca89e;
    font-size: 0.9rem;
    line-height: 1.5;
  }

  .card-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: auto;
    padding-top: 1rem;
    border-top: 1px solid rgba(58, 90, 64, 0.2);
  }

  .price {
    font-size: 1.5rem;
    font-weight: 700;
    color: #5a8a62;
  }

  .btn-add {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.75rem 1.25rem;
    background: transparent;
    border: 1px solid #5a8a62;
    color: #5a8a62;
    border-radius: 8px;
    font-weight: 600;
    font-size: 0.875rem;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-add:hover {
    background: #5a8a62;
    color: #0d1210;
  }

  .btn-add svg {
    width: 16px;
    height: 16px;
  }

  /* Features Section */
  .features-section {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 2rem;
    padding: 4rem 2rem;
    max-width: 1280px;
    margin: 0 auto;
    border-top: 1px solid rgba(58, 90, 64, 0.2);
  }

  @media (max-width: 768px) {
    .features-section {
      grid-template-columns: 1fr;
    }
  }

  .feature {
    text-align: center;
    padding: 2rem;
  }

  .feature-icon {
    width: 56px;
    height: 56px;
    margin: 0 auto 1rem;
    background: rgba(90, 138, 98, 0.15);
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .feature-icon svg {
    width: 28px;
    height: 28px;
    color: #5a8a62;
  }

  .feature h3 {
    font-size: 1.125rem;
    font-weight: 600;
    margin-bottom: 0.5rem;
  }

  .feature p {
    color: #9ca89e;
    font-size: 0.9rem;
  }

  /* Cart Overlay */
  .cart-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.6);
    z-index: 200;
    backdrop-filter: blur(4px);
  }

  /* Cart Sidebar */
  .cart-sidebar {
    position: fixed;
    top: 0;
    right: 0;
    bottom: 0;
    width: 400px;
    max-width: 100%;
    background: #141a17;
    border-left: 1px solid rgba(58, 90, 64, 0.3);
    z-index: 300;
    display: flex;
    flex-direction: column;
    transform: translateX(100%);
    transition: transform 0.3s ease;
  }

  .cart-sidebar.open {
    transform: translateX(0);
  }

  .cart-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem;
    border-bottom: 1px solid rgba(58, 90, 64, 0.2);
  }

  .cart-header h2 {
    font-family: 'Playfair Display', serif;
    font-size: 1.5rem;
    font-weight: 600;
  }

  .close-btn {
    background: transparent;
    border: none;
    color: #9ca89e;
    padding: 0.5rem;
    cursor: pointer;
    border-radius: 8px;
    transition: all 0.2s;
  }

  .close-btn:hover {
    background: rgba(90, 138, 98, 0.15);
    color: #e8ede9;
  }

  .close-btn svg {
    width: 24px;
    height: 24px;
  }

  .empty-cart {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 2rem;
    text-align: center;
    color: #6b736c;
  }

  .empty-icon {
    width: 64px;
    height: 64px;
    margin-bottom: 1rem;
    opacity: 0.5;
  }

  .empty-cart p {
    font-size: 1.125rem;
    color: #9ca89e;
    margin-bottom: 0.25rem;
  }

  .empty-cart span {
    font-size: 0.875rem;
  }

  .cart-items {
    flex: 1;
    overflow-y: auto;
    padding: 1rem 1.5rem;
  }

  .cart-item {
    display: flex;
    gap: 1rem;
    padding: 1rem 0;
    border-bottom: 1px solid rgba(58, 90, 64, 0.15);
  }

  .cart-item-img {
    width: 72px;
    height: 72px;
    object-fit: cover;
    border-radius: 8px;
    background: #1a231e;
  }

  .cart-item-details {
    flex: 1;
  }

  .cart-item-details h4 {
    font-size: 0.95rem;
    font-weight: 600;
    margin-bottom: 0.25rem;
  }

  .cart-item-price {
    color: #5a8a62;
    font-weight: 600;
    font-size: 0.9rem;
    margin-bottom: 0.5rem;
  }

  .quantity-controls {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .quantity-controls button {
    width: 28px;
    height: 28px;
    background: rgba(90, 138, 98, 0.15);
    border: none;
    border-radius: 6px;
    color: #e8ede9;
    cursor: pointer;
    font-size: 1rem;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
  }

  .quantity-controls button:hover {
    background: rgba(90, 138, 98, 0.3);
  }

  .quantity-controls span {
    font-weight: 600;
    min-width: 20px;
    text-align: center;
  }

  .remove-btn {
    background: transparent;
    border: none;
    color: #6b736c;
    padding: 0.5rem;
    cursor: pointer;
    transition: color 0.2s;
  }

  .remove-btn:hover {
    color: #f75a68;
  }

  .remove-btn svg {
    width: 18px;
    height: 18px;
  }

  .cart-footer {
    padding: 1.5rem;
    border-top: 1px solid rgba(58, 90, 64, 0.2);
    background: #0d1210;
  }

  .cart-total {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }

  .cart-total span:first-child {
    color: #9ca89e;
  }

  .total-amount {
    font-size: 1.5rem;
    font-weight: 700;
    color: #e8ede9;
  }

  .btn-checkout {
    width: 100%;
    padding: 1rem;
    background: #5a8a62;
    border: none;
    color: #0d1210;
    font-weight: 700;
    font-size: 1rem;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-checkout:hover {
    background: #4a7a52;
  }

  .shipping-note {
    text-align: center;
    color: #6b736c;
    font-size: 0.8rem;
    margin-top: 0.75rem;
  }

  /* Responsive */
  @media (max-width: 768px) {
    .nav-links {
      display: none;
    }

    .hero {
      padding: 8rem 1.5rem 4rem;
    }

    .hero-title {
      font-size: 2.5rem;
    }

    .catalog-section {
      padding: 3rem 1.5rem;
    }
  }
</style>
