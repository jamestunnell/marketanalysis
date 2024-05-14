import van from "vanjs-core"

const {a} = van.tags

const DownloadJSON = ({name, obj}) => {
    const d = JSON.stringify(obj, null, 2)
    const blob = new Blob([d], {type: 'application/json; charset=utf-8;'})
    const url = URL.createObjectURL(blob);
    const filename = name + ".json";
    const link = a({href: url, download: filename})

    van.add(document.body, link);

    link.click();

    link.remove();

    URL.revokeObjectURL(url);
}

export {DownloadJSON};