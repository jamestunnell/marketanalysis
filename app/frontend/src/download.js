import van from "vanjs-core"

const {a} = van.tags

const Download = ({filename, blob}) => {
    const url = URL.createObjectURL(blob);
    const link = a({href: url, download: filename})

    van.add(document.body, link);

    link.click();

    link.remove();

    URL.revokeObjectURL(url);
};

const DownloadJSON = ({basename, obj}) => {
    const jsonStr = JSON.stringify(obj, null, 2)

    console.log(`downloading JSON ${basename}`)

    Download({
        filename: basename + ".json",
        blob: new Blob([jsonStr], {type: 'application/json; charset=utf-8;'}),
    })
};

export {Download, DownloadJSON};