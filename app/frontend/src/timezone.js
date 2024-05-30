export default function userTimeZone() {
    return Intl.DateTimeFormat().resolvedOptions().timeZone
}