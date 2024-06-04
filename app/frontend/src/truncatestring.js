function truncateString(id, len) {
    if (id.length > len) {
        return id.substring(0, len)
    }
    
    return id
}

function truncateStringAddElipses(id, len) {
    if (id.length > len) {
        return id.substring(0, len) + "..."
    }
    
    return id
}

export {truncateString, truncateStringAddElipses}