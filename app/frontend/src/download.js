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

const DownloadCSV = ({name, data}) => {
    Download({
        filename: name + ".csv",
        blob: new Blob([data], {type: 'text/csv'}),
    })
};

const DownloadJSON = ({name, obj}) => {
    const jsonStr = JSON.stringify(obj, null, 2)

    Download({
        filename: name + ".json",
        blob: new Blob([jsonStr], {type: 'application/json; charset=utf-8;'}),
    })
};

export {Download, DownloadCSV, DownloadJSON};