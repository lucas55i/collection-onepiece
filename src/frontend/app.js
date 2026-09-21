// Carrega utilitários: browser usa <script src="utils.js"> carregado antes,
// Node.js (testes) usa require para acessar as funções puras.
if (typeof require !== 'undefined' && typeof module !== 'undefined') {
  const { toMonthKey, toDateKey, formatMonthLabel, formatDayLabel, formatAcquiredAt } = require('./utils.js');
  global.toMonthKey       = toMonthKey;
  global.toDateKey        = toDateKey;
  global.formatMonthLabel = formatMonthLabel;
  global.formatDayLabel   = formatDayLabel;
  global.formatAcquiredAt = formatAcquiredAt;
}

const API_URL = 'http://localhost:8080/api';

const grid           = document.getElementById('volumes-grid');
const loading        = document.getElementById('loading');
const errorEl        = document.getElementById('error');
const btnSync        = document.getElementById('btn-sync');
const progressText   = document.getElementById('progress-text');
const monthNav       = document.getElementById('month-nav');
const dayNavWrapper  = document.getElementById('day-nav-wrapper');
const dayNav         = document.getElementById('day-nav');
const dayNavLabel    = document.getElementById('day-nav-label');
const btnClearFilter = document.getElementById('btn-clear-filter');

// Estado global
let allVolumes        = [];   // todos os volumes vindos da API
let activeMonthKey    = null; // "YYYY-MM" selecionado
let activeDayKey      = null; // "YYYY-MM-DD" selecionado

// ─── Inicialização ────────────────────────────────────────────────────────────

document.addEventListener('DOMContentLoaded', fetchVolumes);
btnSync.addEventListener('click', syncVolumes);
btnClearFilter.addEventListener('click', clearFilter);

// ─── Busca / Sync ────────────────────────────────────────────────────────────

async function fetchVolumes() {
  showLoading(true);
  hideError();

  try {
    const res = await fetch(`${API_URL}/volumes`);
    if (!res.ok) throw new Error(`Erro ao buscar volumes: ${res.status}`);

    allVolumes = await res.json();

    if (allVolumes.length === 0) {
      loading.textContent = 'Nenhum volume encontrado. Clique em "Sincronizar volumes" para importar.';
      showLoading(true);
      return;
    }

    buildTimeline();
    applyFilter();
    updateProgress();
  } catch (err) {
    showError(err.message);
  } finally {
    showLoading(false);
  }
}

async function syncVolumes() {
  btnSync.disabled = true;
  btnSync.textContent = '⏳ Sincronizando...';
  hideError();

  try {
    const res = await fetch(`${API_URL}/sync`, { method: 'POST' });
    if (!res.ok) throw new Error(`Erro na sincronização: ${res.status}`);

    await fetchVolumes();
  } catch (err) {
    showError(err.message);
  } finally {
    btnSync.disabled = false;
    btnSync.textContent = '🔄 Sincronizar volumes';
  }
}

// ─── Toggle collected ─────────────────────────────────────────────────────────

/**
 * Envia PATCH para o backend e atualiza o volume no estado local.
 * Quando collected=true, o backend define acquired_at = now().
 * Quando collected=false, o backend limpa acquired_at.
 */
async function toggleCollected(id, collected, card, label) {
  try {
    const res = await fetch(`${API_URL}/volumes/${id}`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ collected }),
    });

    if (!res.ok) throw new Error(`Erro ao atualizar: ${res.status}`);

    const updated = await res.json();

    // Atualiza o estado local
    const idx = allVolumes.findIndex(v => v.id === id);
    if (idx !== -1) allVolumes[idx] = updated;

    // Re-renderiza o card com a data atualizada
    renderCard(updated, card);

    // Reconstrói timeline e progresso
    buildTimeline();
    updateProgress();

    // Requirement 3.9: se o período filtrado ficou vazio, limpa o filtro automaticamente
    if (activeMonthKey || activeDayKey) {
      const stillHasVolumes = allVolumes.some(v => {
        if (!v.acquired_at) return false;
        if (activeDayKey) return toDateKey(v.acquired_at) === activeDayKey;
        return toMonthKey(v.acquired_at) === activeMonthKey;
      });
      if (!stillHasVolumes) {
        clearFilter();
        return;
      }
      applyFilter();
    }
  } catch (err) {
    showError(`Falha ao atualizar volume: ${err.message}`);
    // Reverte o checkbox visualmente
    const checkbox = card.querySelector('input[type="checkbox"]');
    if (checkbox) checkbox.checked = !collected;
  }
}

// ─── Renderização ────────────────────────────────────────────────────────────

function applyFilter() {
  let visible = allVolumes;

  if (activeDayKey) {
    visible = allVolumes.filter(v => v.acquired_at && toDateKey(v.acquired_at) === activeDayKey);
  } else if (activeMonthKey) {
    visible = allVolumes.filter(v => v.acquired_at && toMonthKey(v.acquired_at) === activeMonthKey);
  }

  renderVolumes(visible);
}

function renderVolumes(volumes) {
  grid.innerHTML = '';
  volumes.forEach(v => {
    const card = document.createElement('div');
    renderCard(v, card);
    grid.appendChild(card);
  });
}

/**
 * Preenche (ou re-preenche) o conteúdo de um card de volume.
 * Recebe o elemento <div> já existente para poder re-usar em atualizações.
 */
function renderCard(volume, card) {
  card.innerHTML = '';
  card.className = 'volume-card';
  if (volume.collected) card.classList.add('collected');
  card.dataset.id = volume.id;

  // Capa
  const img = document.createElement('img');
  img.src     = volume.cover_image || `https://via.placeholder.com/120x170?text=Vol+${volume.volume_number}`;
  img.alt     = `Capa do Volume ${volume.volume_number}`;
  img.loading = 'lazy';

  // Número
  const number = document.createElement('span');
  number.className   = 'volume-number';
  number.textContent = `Vol. ${volume.volume_number}`;

  // Título
  const title = document.createElement('span');
  title.className   = 'volume-title';
  title.textContent = volume.title || `Volume ${volume.volume_number}`;

  // Data de aquisição
  const dateEl = document.createElement('span');
  dateEl.className   = 'volume-acquired';
  dateEl.textContent = formatAcquiredAt(volume.acquired_at);

  // Checkbox
  const checkboxWrapper = document.createElement('div');
  checkboxWrapper.className = 'checkbox-wrapper';

  const checkbox = document.createElement('input');
  checkbox.type      = 'checkbox';
  checkbox.checked   = volume.collected;
  checkbox.id        = `vol-${volume.id}`;
  checkbox.setAttribute('aria-label', `Marcar volume ${volume.volume_number} como coletado`);

  const label = document.createElement('label');
  label.className   = 'checkbox-label';
  label.htmlFor     = `vol-${volume.id}`;
  label.textContent = volume.collected ? 'Tenho' : 'Não tenho';

  checkbox.addEventListener('change', () =>
    toggleCollected(volume.id, checkbox.checked, card, label)
  );

  // Clique no card alterna o checkbox (exceto quando clica no próprio checkbox/label)
  card.addEventListener('click', (e) => {
    if (e.target !== checkbox && e.target !== label) {
      checkbox.checked = !checkbox.checked;
      toggleCollected(volume.id, checkbox.checked, card, label);
    }
  });

  checkboxWrapper.appendChild(checkbox);
  checkboxWrapper.appendChild(label);

  card.appendChild(img);
  card.appendChild(number);
  card.appendChild(title);
  card.appendChild(dateEl);
  card.appendChild(checkboxWrapper);
}

// ─── Timeline ────────────────────────────────────────────────────────────────

/**
 * Constrói a navegação de meses e dias a partir dos volumes adquiridos.
 */
function buildTimeline() {
  // Agrupa volumes por mês e por dia
  const byMonth = {}; // { "2024-01": { label, days: { "2024-01-15": [volumes] } } }

  allVolumes.forEach(v => {
    if (!v.acquired_at) return;

    const mk = toMonthKey(v.acquired_at);
    const dk = toDateKey(v.acquired_at);

    if (!byMonth[mk]) {
      byMonth[mk] = { label: formatMonthLabel(v.acquired_at), days: {} };
    }
    if (!byMonth[mk].days[dk]) {
      byMonth[mk].days[dk] = { label: formatDayLabel(v.acquired_at), count: 0 };
    }
    byMonth[mk].days[dk].count++;
  });

  const monthKeys = Object.keys(byMonth).sort();

  // Esconde a timeline se não houver nenhum adquirido
  const timelineNav = document.getElementById('timeline-nav');
  if (monthKeys.length === 0) {
    timelineNav.classList.add('hidden');
    return;
  }
  timelineNav.classList.remove('hidden');

  // Renderiza os botões de mês
  monthNav.innerHTML = '';
  monthKeys.forEach(mk => {
    const totalDia = Object.values(byMonth[mk].days).reduce((s, d) => s + d.count, 0);
    const btn = document.createElement('button');
    btn.className     = 'timeline-btn month-btn' + (mk === activeMonthKey ? ' active' : '');
    btn.dataset.month = mk;
    btn.setAttribute('role', 'listitem');
    btn.setAttribute('aria-pressed', mk === activeMonthKey ? 'true' : 'false');
    btn.innerHTML     = `<span class="tl-label">${byMonth[mk].label}</span><span class="tl-badge">${totalDia}</span>`;

    btn.addEventListener('click', () => {
      if (activeMonthKey === mk) {
        // Segundo clique no mesmo mês: remove o filtro de mês
        clearFilter();
      } else {
        activeMonthKey = mk;
        activeDayKey   = null;
        buildTimeline();     // re-renderiza highlight
        buildDayNav(byMonth[mk]);
        applyFilter();
        updateClearBtn();
      }
    });

    monthNav.appendChild(btn);
  });

  // Se um mês ainda está ativo, reconstrói o nav de dias
  if (activeMonthKey && byMonth[activeMonthKey]) {
    buildDayNav(byMonth[activeMonthKey]);
  } else {
    dayNavWrapper.classList.add('hidden');
  }
}

function buildDayNav(monthData) {
  dayNavWrapper.classList.remove('hidden');
  dayNavLabel.textContent = monthData.label;
  dayNav.innerHTML = '';

  const dayKeys = Object.keys(monthData.days).sort();
  dayKeys.forEach(dk => {
    const d = monthData.days[dk];
    const btn = document.createElement('button');
    btn.className    = 'timeline-btn day-btn' + (dk === activeDayKey ? ' active' : '');
    btn.dataset.day  = dk;
    btn.setAttribute('role', 'listitem');
    btn.setAttribute('aria-pressed', dk === activeDayKey ? 'true' : 'false');
    btn.innerHTML    = `<span class="tl-label">${d.label}</span><span class="tl-badge">${d.count}</span>`;

    btn.addEventListener('click', () => {
      if (activeDayKey === dk) {
        // Segundo clique: volta para filtro só por mês
        activeDayKey = null;
        buildTimeline();
        applyFilter();
        updateClearBtn();
      } else {
        activeDayKey = dk;
        buildTimeline();
        applyFilter();
        updateClearBtn();
      }
    });

    dayNav.appendChild(btn);
  });
}

function clearFilter() {
  activeMonthKey = null;
  activeDayKey   = null;
  buildTimeline();
  applyFilter();
  updateClearBtn();
}

function updateClearBtn() {
  if (activeMonthKey || activeDayKey) {
    btnClearFilter.classList.remove('hidden');
  } else {
    btnClearFilter.classList.add('hidden');
  }
}

// ─── Progresso ───────────────────────────────────────────────────────────────

function updateProgress() {
  const total     = allVolumes.length;
  const collected = allVolumes.filter(v => v.collected).length;
  progressText.textContent = `${collected} / ${total} volumes coletados`;
}

// ─── UI helpers ──────────────────────────────────────────────────────────────

function showLoading(show) {
  loading.classList.toggle('hidden', !show);
}

function showError(msg) {
  errorEl.textContent = msg;
  errorEl.classList.remove('hidden');
}

function hideError() {
  errorEl.classList.add('hidden');
}
