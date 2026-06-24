document.addEventListener('DOMContentLoaded', () => {
    const section = document.querySelector('.recent-section');
    const track = document.getElementById('recent-watched');
    if (!section || !track) {
        return;
    }

    function render() {
        const history = getHistory();
        track.innerHTML = '';

        if (history.length === 0) {
            section.hidden = true;
            return;
        }

        section.hidden = false;

        history.forEach(item => {
            const href = item.type === 'tv'
                ? `/tv/${encodeURIComponent(item.id)}?s=${encodeURIComponent(item.season || 1)}&e=${encodeURIComponent(item.episode || 1)}&srv=${encodeURIComponent(item.server || '')}`
                : `/movie/${encodeURIComponent(item.id)}?srv=${encodeURIComponent(item.server || '')}`;

            const card = document.createElement('article');
            card.className = 'recent-card';
            card.dataset.id = item.id;

            const metaBits = [];
            metaBits.push(`<span class="recent-card__type recent-card__type--${item.type}">${item.type === 'tv' ? 'TV' : 'Movie'}</span>`);
            if (item.type === 'tv') {
                metaBits.push(`<span class="recent-card__meta">S${escapeHtml(item.season || 1)} · E${escapeHtml(item.episode || 1)}</span>`);
            }

            card.innerHTML = `
                <button class="recent-card__remove" type="button" aria-label="Remove from history" title="Remove">×</button>
                <a class="recent-card__link" href="${href}" aria-label="${escapeHtml(item.title)}"></a>
                <div class="recent-card__poster">
                    <img src="${escapeHtml(item.poster)}" alt="${escapeHtml(item.title)}" loading="lazy">
                </div>
                <div class="recent-card__body">
                    <h3 class="recent-card__title">${escapeHtml(item.title)}</h3>
                    <div class="recent-card__footer">${metaBits.join('')}</div>
                </div>
            `;

            card.querySelector('.recent-card__remove').addEventListener('click', e => {
                e.preventDefault();
                e.stopPropagation();
                removeLastVisited(item.id);
                render();
                updateNav();
            });

            track.appendChild(card);
        });

        updateNav();
    }

    function updateNav() {
        const prev = section.querySelector('.recent-nav--prev');
        const next = section.querySelector('.recent-nav--next');
        if (!prev || !next) {
            return;
        }
        const maxScroll = track.scrollWidth - track.clientWidth;
        prev.disabled = track.scrollLeft <= 0;
        next.disabled = track.scrollLeft >= maxScroll - 1;
        const overflow = maxScroll > 4;
        prev.hidden = !overflow;
        next.hidden = !overflow;
    }

    const prevBtn = section.querySelector('.recent-nav--prev');
    const nextBtn = section.querySelector('.recent-nav--next');
    if (prevBtn) {
        prevBtn.addEventListener('click', () => {
            track.scrollBy({left: -track.clientWidth * 0.8, behavior: 'smooth'});
        });
    }
    if (nextBtn) {
        nextBtn.addEventListener('click', () => {
            track.scrollBy({left: track.clientWidth * 0.8, behavior: 'smooth'});
        });
    }

    track.addEventListener('scroll', updateNav);
    window.addEventListener('resize', updateNav);

    render();
});

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text ?? '';
    return div.innerHTML;
}
