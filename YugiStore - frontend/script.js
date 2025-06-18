// ===================================================
// ================  CONFIGURAÇÕES GERAIS  ===========
// ===================================================

// const API_URL = "http://localhost:8080";
const API_URL = "https://yugistore-650237781960.us-central1.run.app"

// ===================================================
// ====================== LOGIN ======================
// ===================================================

import { initializeApp } from "https://www.gstatic.com/firebasejs/11.7.3/firebase-app.js";
import {
    getAuth,
    GoogleAuthProvider,
    signInWithPopup,
    signOut,
    onAuthStateChanged,
    setPersistence,
    browserLocalPersistence
} from "https://www.gstatic.com/firebasejs/11.7.3/firebase-auth.js";

// Firebase app configuration
const firebaseConfig = {
    apiKey: "AIzaSyAu7fmA4io9jDTxnVXvgmGT1Z6E8VhDeIM",
    authDomain: "yugistore.firebaseapp.com",
    projectId: "yugistore",
    storageBucket: "yugistore.firebasestorage.app",
    messagingSenderId: "650237781960",
    appId: "1:650237781960:web:f620e478c591cb86a602fd"
};

const app = initializeApp(firebaseConfig);
const auth = getAuth(app);
const provider = new GoogleAuthProvider();
export let token, user, credential;

setPersistence(auth, browserLocalPersistence)
    .catch(err => console.error("Erro ao definir persistência:", err));

const loginBtn = document.getElementById("loginBtn");
const authArea = document.getElementById("authArea");

onAuthStateChanged(auth, u => {
    if (u) {
        user = u;
        renderUserUI(u);
    }
});

// Render username dropdown
function renderUserUI(u) {
    authArea.innerHTML = `
    <div class="dropdown">
      <button class="btn btn-outline-light dropdown-toggle"
              id="userDropdown"
              data-bs-toggle="dropdown"
              aria-expanded="false">
        ${u.displayName || "Usuário"}
      </button>
      <ul class="dropdown-menu dropdown-menu-end">
        <li><a class="dropdown-item" href="#" id="logoutBtn">Logout</a></li>
      </ul>
    </div>`;
    document.getElementById("logoutBtn").addEventListener("click", e => {
        e.preventDefault();
        signOut(auth).then(() => location.reload());
    });
}

// Login click
loginBtn.addEventListener("click", () => {
    signInWithPopup(auth, provider)
        .then(result => {
            credential = GoogleAuthProvider.credentialFromResult(result);
            token = credential.accessToken;
            user = result.user;
            renderUserUI(user);
            console.log("TOKEN:", token);
            console.log("USER :", user);
        })
        .catch(err => console.error(err.code, err.message));
});

// ====================================================
// ===================== CARRINHO =====================
// ====================================================

const cartSidebar = document.getElementById("cartSidebar");
const cartPullTab = document.getElementById("cartPullTab");
const closeCartBtn = document.getElementById("closeCartBtn");
const cartItemsList = document.getElementById("cartItems");
const cartTotal = document.getElementById("cartTotal");
const checkoutBtn = document.getElementById("checkoutBtn");

let cart = JSON.parse(localStorage.getItem("cart")) || [];

function saveCart() {
    localStorage.setItem("cart", JSON.stringify(cart));
}

function openCart() {
    cartSidebar.classList.add("open");
}
function closeCart() {
    cartSidebar.classList.remove("open");
}
cartPullTab.addEventListener("click", () => {
    cartSidebar.classList.toggle("open");
});
closeCartBtn.addEventListener("click", closeCart);

function updateCartUI() {
    cartItemsList.innerHTML = "";
    let total = 0;
    if (cart.length === 0) {
        cartItemsList.innerHTML = "<li class='list-group-item text-muted'>Empty Cart.</li>";
        checkoutBtn.disabled = true;
    } else {
        cart.forEach((item, idx) => {
            total += item.price * item.qty;
            const li = document.createElement("li");
            li.className = "list-group-item d-flex align-items-center justify-content-between";
            li.innerHTML = `
                <img src="${item.image}" alt="${item.name}" width="36" height="36" class="me-2 rounded" style="object-fit:cover;">
                <span class="flex-grow-1">${item.name} <span class="badge bg-secondary ms-1">x${item.qty}</span></span>
                <span>${(item.price * item.qty).toFixed(2)} R$</span>
                <button class="btn btn-sm btn-outline-danger ms-2" title="Remover" data-idx="${idx}">&times;</button>
            `;
            li.querySelector("button").onclick = () => {
                // Diminui quantidade, remove se qty==1
                if (item.qty > 1) {
                    item.qty--;
                } else {
                    cart.splice(idx, 1);
                }
                saveCart();
                updateCartUI();
            };
            cartItemsList.appendChild(li);
        });
        checkoutBtn.disabled = false;
    }
    cartTotal.textContent = total.toFixed(2) + " R$";
}

updateCartUI();

// Adicione isso na função que já existe no seu código:
function addToCart(card) {
    const existing = cart.find(item =>
        item.edition === card.edition && item.state === card.state
    );

    if (existing) {
        existing.qty += 1;
    } else {
        cart.push({ ...card, qty: 1 });
    }
    saveCart();
    updateCartUI();
    openCart();
}



// Garantir persistência mesmo se recarregar a página
window.addEventListener("DOMContentLoaded", updateCartUI);

checkoutBtn.addEventListener("click", () => {
    alert(":D");
});


// ===================================================
// ====================== BUSCA ======================
// ===================================================

document.addEventListener("DOMContentLoaded", () => {
    document.getElementById("searchCardForm").addEventListener("submit", async (e) => {
        e.preventDefault();
        const name = document.getElementById("searchName").value.trim();
        let url = `${API_URL}/card`;
        if (name) {
            url += `?name=${encodeURIComponent(name)}`;
        }

        try {
            const cardRes = await fetch(url);
            const cards = await cardRes.json();

            const allResults = [];

            for (const card of cards) {
                const prodRes = await fetch(`${API_URL}/product?edition=${encodeURIComponent(card.edition)}`);
                if (!prodRes.ok) continue;

                const products = await prodRes.json();
                if (!Array.isArray(products)) continue;

                for (const product of products) {
                    allResults.push({
                        name: card.name,
                        edition: card.edition,
                        image_url: card.image_url,
                        state: product.state,
                        price: parseFloat(product.price)
                    });
                }
            }

            renderSearchResults(allResults);
        } catch (err) {
            console.error("Error searching:", err);
        }
    });

    // Executar busca inicial automática
    document.getElementById("searchCardForm").dispatchEvent(new Event("submit"));
});


function renderSearchResults(cards) {
    const list = document.getElementById("searchResults");
    list.innerHTML = "";

    if (!cards.length) {
        list.innerHTML = '<li class="text-muted" style="color: white !important;">No card found.</li>';
        return;
    }

    cards.forEach(card => {
        const li = document.createElement("li");
        li.className = "col-6 col-sm-4 col-md-3 col-lg-2";
        li.innerHTML = `
        <div class="card h-100">
            <div class="card-img-container">
                <img src="${card.image_url}" alt="${card.name}">
            </div>
            <div class="card-body d-flex flex-column">
                <h6 class="card-title mb-1 card-name">${card.name}</h6>
                <p class="card-text mb-1 card-edition">Edition: ${card.edition}</p>
                <p class="card-text mb-1 card-state">State: ${card.state}</p>
                <p class="card-text mb-2 fw-bold text-success card-price">${card.price.toFixed(2)} R$</p>
                <button class="btn btn-primary mt-auto">Add to cart</button>
            </div>
        </div>
        `;

        li.querySelector("button").onclick = () =>
            addToCart({
                name: card.name,
                edition: card.edition,
                price: card.price,
                image: card.image_url,
                state: card.state
            });

        list.appendChild(li);
    });
}


// =================================================
// ================= REGISTER CARD =================
// =================================================
document.getElementById("regCardForm").addEventListener("submit", e => {
    e.preventDefault();
    const form = e.target;
    const data = new FormData(form);

    document.getElementById("regCardResult").textContent = "Sending...";

    fetch(`${API_URL}/card`, { method: "POST", body: data })
        .then(res => res.json()).then(res => {
            document.getElementById("regCardResult").textContent = res.status || "Registered card.";
            form.reset();
        })
        .catch(err => {
            document.getElementById("regCardResult").textContent = "Error registering.";
            console.error(err);
        });
});

// ===================================================
// ================ REGISTRAR ESTOQUE ================
// ===================================================

const annEditionSelect = document.getElementById("annEdition");

const fetchAndFillEditions = async () => {
    annEditionSelect.innerHTML = `<option value="">Loading...</option>`;

    try {
        const res = await fetch(`${API_URL}/card`);
        const cards = await res.json();

        if (!Array.isArray(cards)) throw new Error("Formato inválido");

        if (cards.length === 0) {
            annEditionSelect.innerHTML = `<option value="">No edition found</option>`;
            return;
        }

        annEditionSelect.innerHTML = cards.map(code =>
            `<option value="${code.edition}">${code.name} - ${code.edition}</option>`
        ).join("");

    } catch (err) {
        console.error("Erro loading editions:", err);
        annEditionSelect.innerHTML = `<option value="">Erro loading editions</option>`;
    }
};

const announceModal = document.getElementById('announceCardModal');

announceModal.addEventListener('show.bs.modal', () => {
    fetchAndFillEditions();
});

document.getElementById("annCardForm").addEventListener("submit", function (e) {
    e.preventDefault();

    const edition = document.getElementById("annEdition").value;
    const state = document.getElementById("annState").value;
    const price = document.getElementById("annPrice").value;

    const data = {
        edition: edition,
        state: state,
        price: price
    };

    console.log(data)

    fetch(`${API_URL}/product`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    })
        .then(res => res.json())
        .then(res => {
            document.getElementById("annCardResult").textContent = res.status || "Card successfully announced!";
            document.getElementById("annCardForm").reset();
        })
        .catch(err => {
            console.error("Error announcing card:", err);
            document.getElementById("annCardResult").textContent = "Error announcing card.";
        });
});

