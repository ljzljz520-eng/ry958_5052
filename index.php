<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>校园文创商品站</title>
  <link rel="stylesheet" href="/styles.css">
</head>
<body>
  <header class="site-header">
    <div class="shell header-inner">
      <a class="brand" href="/">校园文创</a>
      <div class="member-status" id="member-status">未登录</div>
    </div>
  </header>
  <main class="shell">
    <section class="intro">
      <p class="eyebrow">CAMPUS CREATIVE GOODS</p>
      <h1>把校园记忆带回日常</h1>
      <p class="intro-copy">帆布袋、徽章、明信片和笔记本，都是为校园生活准备的小小纪念。</p>
    </section>
    <section class="catalog-section" aria-labelledby="catalog-title">
      <div class="section-heading">
        <div>
          <p class="eyebrow">精选商品</p>
          <h2 id="catalog-title">校园日用与纪念</h2>
        </div>
        <span class="section-note">固定夹具商品</span>
      </div>
      <div class="product-grid" id="product-grid">
        <article class="product-card" data-product="canvas-bag"><div class="product-art art-bag">帆布袋</div><h3>校园帆布袋</h3><p>轻便耐用的校园日常帆布袋</p><strong>¥39.00</strong><button class="add-button" data-product-id="canvas-bag">加入购物车</button></article>
        <article class="product-card" data-product="campus-badge"><div class="product-art art-badge">徽章</div><h3>校园徽章</h3><p>别在书包上的校园纪念徽章</p><strong>¥12.00</strong><button class="add-button" data-product-id="campus-badge">加入购物车</button></article>
        <article class="product-card" data-product="postcard"><div class="product-art art-postcard">明信片</div><h3>校园明信片</h3><p>记录校园风景的明信片套装</p><strong>¥16.00</strong><button class="add-button" data-product-id="postcard">加入购物车</button></article>
        <article class="product-card" data-product="notebook"><div class="product-art art-notebook">笔记本</div><h3>校园笔记本</h3><p>适合课程记录的线圈笔记本</p><strong>¥22.00</strong><button class="add-button" data-product-id="notebook">加入购物车</button></article>
      </div>
    </section>
    <section class="account-section" aria-labelledby="account-title">
      <div class="account-panel">
        <p class="eyebrow">会员中心</p>
        <h2 id="account-title">登录后管理购物车</h2>
        <div class="forms">
          <form id="register-form"><h3>注册会员</h3><label>用户名<input name="username" required autocomplete="username"></label><label>密码<input name="password" type="password" required autocomplete="new-password"></label><button type="submit">注册</button></form>
          <form id="login-form"><h3>会员登录</h3><label>用户名<input name="username" required autocomplete="username"></label><label>密码<input name="password" type="password" required autocomplete="current-password"></label><button type="submit">登录</button></form>
        </div>
        <p class="feedback" id="feedback" role="status"></p>
      </div>
      <aside class="cart-panel" aria-labelledby="cart-title"><div class="cart-heading"><h2 id="cart-title">我的购物车</h2><span id="cart-total">¥0.00</span></div><ul id="cart-items"><li>登录后查看购物车</li></ul></aside>
    </section>
  </main>
  <script src="/app.js"></script>
</body>
</html>
