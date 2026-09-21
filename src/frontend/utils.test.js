const { describe, it } = require('node:test');
const assert = require('node:assert');
const { toMonthKey, toDateKey, formatAcquiredAt } = require('./utils.js');

describe('toMonthKey', () => {
  it('retorna YYYY-MM a partir de uma string ISO com Z', () => {
    assert.strictEqual(toMonthKey('2024-01-15T14:32:00Z'), '2024-01');
  });
  it('retorna YYYY-MM a partir de uma string ISO com offset positivo', () => {
    assert.strictEqual(toMonthKey('2024-03-20T10:00:00+03:00'), '2024-03');
  });
});

describe('toDateKey', () => {
  it('retorna YYYY-MM-DD a partir de uma string ISO com Z', () => {
    assert.strictEqual(toDateKey('2024-01-15T14:32:00Z'), '2024-01-15');
  });
  it('retorna YYYY-MM-DD a partir de uma string ISO com offset negativo', () => {
    assert.strictEqual(toDateKey('2024-12-31T23:59:59-05:00'), '2024-12-31');
  });
});

describe('formatAcquiredAt', () => {
  it('retorna string vazia para null', () => {
    assert.strictEqual(formatAcquiredAt(null), '');
  });
  it('retorna string vazia para undefined', () => {
    assert.strictEqual(formatAcquiredAt(undefined), '');
  });
  it('retorna string contendo "às" para ISO válido', () => {
    const result = formatAcquiredAt('2024-01-15T14:32:00Z');
    assert.ok(result.includes('às'), `esperado "às" em: ${result}`);
  });
});
