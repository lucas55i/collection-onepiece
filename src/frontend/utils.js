/**
 * utils.js — funções puras extraídas de app.js para permitir testes isolados.
 * Compatível com browser (sem module.exports) e Node.js (com module.exports).
 */

/** Retorna "YYYY-MM" a partir de uma string ISO */
function toMonthKey(isoStr) {
  return isoStr.slice(0, 7);
}

/** Retorna "YYYY-MM-DD" a partir de uma string ISO */
function toDateKey(isoStr) {
  return isoStr.slice(0, 10);
}

/** Formata o label do mês: "janeiro de 2024" */
function formatMonthLabel(isoStr) {
  const date = new Date(isoStr);
  return date.toLocaleDateString('pt-BR', { month: 'long', year: 'numeric' });
}

/** Formata o label do dia: "seg., 15" */
function formatDayLabel(isoStr) {
  const date = new Date(isoStr);
  return date.toLocaleDateString('pt-BR', { day: '2-digit', weekday: 'short' });
}

/** Formata a data exibida no card: "🗓 15/01/2024 às 14:32" ou string vazia */
function formatAcquiredAt(isoStr) {
  if (!isoStr) return '';
  const date = new Date(isoStr);
  const d = date.toLocaleDateString('pt-BR', { day: '2-digit', month: '2-digit', year: 'numeric' });
  const t = date.toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' });
  return `🗓 ${d} às ${t}`;
}

// Compatibilidade Node.js / browser
if (typeof module !== 'undefined') {
  module.exports = { toMonthKey, toDateKey, formatMonthLabel, formatDayLabel, formatAcquiredAt };
}
