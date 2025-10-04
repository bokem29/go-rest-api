 let editId = null;
    let token = localStorage.getItem('authToken') || '';
    let refreshToken = localStorage.getItem('refreshToken') || '';

    function refreshAuthUI() {
      const loginDiv = document.getElementById('loginForm');
      const loggedInDiv = document.getElementById('loggedIn');
      const statusText = document.getElementById('statusText');
      if (token) {
        loginDiv.style.display = 'none';
        loggedInDiv.style.display = 'block';
        statusText.textContent = 'Logged in';
      } else {
        loginDiv.style.display = 'block';
        loggedInDiv.style.display = 'none';
        statusText.textContent = 'Belum login';
      }
    }

    async function login() {
      const username = document.getElementById('username').value;
      const password = document.getElementById('password').value;
      const res = await fetch('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
      });
      if (!res.ok) {
        showToast('Login gagal', 'error');
        return;
      }
      const data = await res.json();
      token = data.token;
      refreshToken = data.refresh;
      localStorage.setItem('authToken', token);
      localStorage.setItem('refreshToken', refreshToken);
      refreshAuthUI();
      showToast('Login berhasil', 'success');
    }

    async function logout() {
      if (!token) return;
      await fetch('/api/logout', { method: 'POST', headers: { 'Authorization': 'Bearer ' + token }});
      token = '';
      refreshToken = '';
      localStorage.removeItem('authToken');
      localStorage.removeItem('refreshToken');
      refreshAuthUI();
      document.getElementById('characterCards').innerHTML = '';
      showToast('Berhasil logout');
    }

    async function apiFetch(url, options = {}, retry = true) {
      if (!options.headers) options.headers = {};
      if (token) options.headers['Authorization'] = 'Bearer ' + token;
      const res = await fetch(url, options);
      if (res.status === 401 && retry && refreshToken) {
        const ok = await tryRefresh();
        if (ok) {
          return apiFetch(url, options, false);
        }
        // refresh gagal: paksa logout agar UI sinkron jika token dicabut dari terminal
        forceLogoutUI();
      }
      return res;
    }

    async function tryRefresh() {
      try {
        const res = await fetch('/api/refresh', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ refresh: refreshToken })
        });
        if (!res.ok) return false;
        const data = await res.json();
        token = data.token;
        refreshToken = data.refresh;
        localStorage.setItem('authToken', token);
        localStorage.setItem('refreshToken', refreshToken);
        refreshAuthUI();
        return true;
      } catch (e) {
        return false;
      }
    }

    function forceLogoutUI() {
      token = '';
      refreshToken = '';
      localStorage.removeItem('authToken');
      localStorage.removeItem('refreshToken');
      refreshAuthUI();
    }

    async function loadCharacters() {
      if (!token) { alert('Silakan login dulu'); return; }
      document.getElementById('loading').style.display = 'flex';
      const res = await apiFetch("/api/characters");
      const data = await res.json();
      document.getElementById('loading').style.display = 'none';

      const container = document.getElementById("characterCards");
      container.innerHTML = "";
      document.getElementById('emptyState').style.display = (!Array.isArray(data) || data.length === 0) ? 'block' : 'none';

      data.forEach(c => {
        const card = document.createElement("div");
        card.className = "card";

        card.innerHTML = `
          <h3>${c.name}</h3>
          <p><strong>Role:</strong> ${c.role}</p>
          <p><strong>Game:</strong> ${c.game}</p>
          <div style="margin-top: 15px; text-align: center;">
            <button class="btn-edit" onclick="editCharacter(${c.id}, '${c.name}', '${c.role}', '${c.game}')">Edit</button>
            <button class="btn-delete" onclick="deleteCharacter(${c.id})">Hapus</button>
          </div>
        `;
        container.appendChild(card);
      });
    }

    function editCharacter(id, name, role, game) {
      editId = id;
      document.getElementById('modalTitle').textContent = 'Edit Karakter';
      document.getElementById("name").value = name;
      document.getElementById("role").value = role;
      document.getElementById("game").value = game;
      openModal(true);
    }


    async function saveCharacter(event) {
  event.preventDefault();
      if (!token) { alert('Silakan login dulu'); return; }

  const character = {
    name: document.getElementById("name").value,
    role: document.getElementById("role").value,
    game: document.getElementById("game").value
  };

  if (editId) {
    await apiFetch(`/api/characters/${editId}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(character)
    });
    editId = null;
    showToast('Perubahan tersimpan', 'success');
  } else {
    await apiFetch("/api/characters", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(character)
    });
    showToast('Karakter ditambahkan', 'success');
  }

  document.querySelector("#modalBackdrop form").reset();
  closeModal();
  loadCharacters();
}


    async function deleteCharacter(id) {
      if (confirm("Yakin mau hapus karakter ini?")) {
        if (!token) { alert('Silakan login dulu'); return; }
        await apiFetch(`/api/characters/${id}`, { method: "DELETE" });
        loadCharacters();
        showToast('Karakter dihapus', 'success');
      }
    }

    function showToast(msg, type = 'success') {
      const box = document.getElementById('toasts');
      if (!box) {
        const div = document.createElement('div');
        div.id = 'toasts';
        div.className = 'toasts';
        document.body.appendChild(div);
      }
      const t = document.createElement('div');
      t.className = `toast ${type}`;
      t.textContent = msg;
      document.getElementById('toasts').appendChild(t);
      setTimeout(() => { t.remove(); }, 2500);
    }

    function openModal(isEdit = false) {
      if (!token) { alert('Silakan login dulu'); return; }
      if (!isEdit) {
        editId = null;
        document.getElementById('modalTitle').textContent = 'Tambah Karakter';
        document.querySelector("#modalBackdrop form").reset();
      }
      document.getElementById('modalBackdrop').style.display = 'flex';
    }

    function closeModal() {
      document.getElementById('modalBackdrop').style.display = 'none';
      editId = null;
    }

    // Close modal on backdrop click
    document.addEventListener('DOMContentLoaded', function() {
      document.getElementById('modalBackdrop').addEventListener('click', function(e) {
        if (e.target === this) {
          closeModal();
        }
      });
    });
    refreshAuthUI();