const encounter = JSON.parse(document.getElementById('encounter-data').textContent).event;
const notes = {
  accept: 'A measured resolution. Completion earns 5 respect. Existing police attention can complicate the job.',
  'approach:careful': 'More time, a smaller payday, no added heat. Completion still earns 5 respect.',
  'approach:press': 'A faster, better-paying approach with added police attention. Completion earns 5 respect; existing heat can complicate the job.',
  decline: 'Leave this arrangement to somebody else. No time, payment or added heat.'
};
const approaches = [...document.querySelectorAll('[data-choice]')];
for (const button of approaches) button.addEventListener('click', () => {
  const choice = encounter.choices.find(choice => choice.id === button.dataset.choice);
  if (!choice) return;
  for (const other of approaches) other.setAttribute('aria-pressed', String(other === button));
  document.getElementById('choice-time').textContent = choice.minutes ? `${choice.minutes} min` : 'None';
  document.getElementById('choice-pay').textContent = choice.pay ? `$${choice.pay}` : 'None';
  document.getElementById('choice-heat').textContent = choice.heat ? `+${choice.heat}` : 'None';
  document.getElementById('choice-note').textContent = notes[choice.id];
});
const gallery = [
  {src:'assets/bellwether-city.jpg',caption:'Bellwether — a persistent city. Development screenshot.'},
  {src:'assets/saint-agnes.jpg',caption:'Saint Agnes — meet the people behind the names. Development screenshot.'},
  {src:'assets/green-baize.jpg',caption:'The Green Baize — the neighborhood pool hall. Development screenshot.'}
];
const dialog = document.getElementById('lightbox');
for (const button of document.querySelectorAll('[data-gallery]')) button.addEventListener('click', () => {
  const photo = gallery[Number(button.dataset.gallery)];
  const image = document.getElementById('lightbox-image');
  image.src = photo.src; image.alt = photo.caption;
  document.getElementById('lightbox-caption').textContent = photo.caption;
  dialog.showModal();
});
document.getElementById('close-lightbox').addEventListener('click', () => dialog.close());
dialog.addEventListener('click', event => { if (event.target === dialog && (event.clientX < dialog.getBoundingClientRect().left || event.clientX > dialog.getBoundingClientRect().right || event.clientY < dialog.getBoundingClientRect().top || event.clientY > dialog.getBoundingClientRect().bottom)) dialog.close(); });
