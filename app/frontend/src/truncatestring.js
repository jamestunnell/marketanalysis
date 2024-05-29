export default function truncateString(id, len) {
    if (id.length > len) {
        return id.substring(0, len) + "..."
    }
    
    return id
}