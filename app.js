(function () {
  var member = null;
  var feedback = document.getElementById("feedback");
  var status = document.getElementById("member-status");
  var items = document.getElementById("cart-items");
  var total = document.getElementById("cart-total");

  function show(message) { feedback.textContent = message; }
  function request(path, options) {
    return fetch(path, Object.assign({ headers: { "Content-Type": "application/json" } }, options || {})).then(function (response) {
      return response.json().then(function (body) { if (!response.ok) throw new Error(body.error || "请求失败"); return body; });
    });
  }
  function refreshCart() {
    if (!member) return;
    request("/api/cart?member_id=" + encodeURIComponent(member.id), { headers: {} }).then(function (cart) {
      items.innerHTML = cart.items.length ? cart.items.map(function (item) { return "<li>" + item.product.name + " × " + item.quantity + " · ¥" + item.subtotal + "</li>"; }).join("") : "<li>购物车还是空的</li>";
      total.textContent = "¥" + cart.total;
    }).catch(function (error) { show(error.message); });
  }
  function setMember(value) { member = value; status.textContent = "你好，" + value.username; show("已登录，欢迎回来"); refreshCart(); }
  document.getElementById("register-form").addEventListener("submit", function (event) {
    event.preventDefault();
    var form = new FormData(event.target);
    request("/api/register", { method: "POST", body: JSON.stringify({ username: form.get("username"), password: form.get("password") }) }).then(function (body) { setMember(body.member); }).catch(function (error) { show(error.message); });
  });
  document.getElementById("login-form").addEventListener("submit", function (event) {
    event.preventDefault();
    var form = new FormData(event.target);
    request("/api/login", { method: "POST", body: JSON.stringify({ username: form.get("username"), password: form.get("password") }) }).then(function (body) { setMember(body.member); }).catch(function (error) { show(error.message); });
  });
  document.querySelectorAll(".add-button").forEach(function (button) {
    button.addEventListener("click", function () {
      if (!member) { show("请先注册或登录会员"); return; }
      request("/api/cart/items", { method: "POST", body: JSON.stringify({ member_id: member.id, product_id: button.dataset.productId, quantity: 1 }) }).then(function () { show("商品已加入购物车"); refreshCart(); }).catch(function (error) { show(error.message); });
    });
  });
}());
