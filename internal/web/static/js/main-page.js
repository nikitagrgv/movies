document.addEventListener('DOMContentLoaded', () => {
    const container = document.getElementById('recent-watched');
    if (!container) {
        return;
    }

    const history = getHistory();
    if (history.length === 0) {
        return;
    }

    container.innerHTML = '';

    history.forEach(item => {
        const href = item.type === 'tv'
            ? `/tv/${encodeURIComponent(item.id)}?s=${encodeURIComponent(item.season || 1)}&e=${encodeURIComponent(item.episode || 1)}&srv=${encodeURIComponent(item.server || '')}`
            : `/movie/${encodeURIComponent(item.id)}?srv=${encodeURIComponent(item.server || '')}`;

        const article = document.createElement('article');
        article.className = 'movie-card';
        article.innerHTML = `
            <a class="movie-card-link" href="${href}" aria-label="${escapeHtml(item.title)}"></a>
            <img class="movie-poster" src="${escapeHtml(item.poster)}" alt="${escapeHtml(item.title)}" loading="lazy">
            <div class="movie-info">
                <h2 class="movie-title">
                    <a href="${href}">${escapeHtml(item.title)}</a>
                </h2>
                <div class="movie-card__footer">
                    <span class="movie-card__type">${item.type === 'tv' ? 'TV Show' : 'Movie'}</span>
                </div>
            </div>
        `;
        container.appendChild(article);
    });

    const section = container.closest('.recent-section');
    if (section) {
        section.hidden = false;
    }
});

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text ?? '';
    return div.innerHTML;
}
