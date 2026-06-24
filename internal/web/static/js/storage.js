const recentWatchedKey = 'recentWatched';

function getHistory() {
    return JSON.parse(localStorage.getItem(recentWatchedKey)) || [];
}

function saveHistory(history) {
    localStorage.setItem(recentWatchedKey, JSON.stringify(history));

}

function addLastVisited(mediaItem) {
    let history = getHistory();

    history = history.filter(item => item.id !== mediaItem.id)

    history.unshift(mediaItem);

    if (history.length > 20) {
        history.length = 20;
    }

    saveHistory(history);
}

function removeLastVisited(id) {
    let history = getHistory();
    history = history.filter(item => String(item.id) !== String(id));
    saveHistory(history);
}
