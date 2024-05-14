import van from "vanjs-core"

const {input} = van.tags

const ReadFileInputAsText = () => {
    return new Promise((resolve, reject) => {
        const fileInput = input({
            type:"file",
            onchange: e => {
                var file = e.target.files[0]
                
                console.log("importing file", file.name)
                
                const reader = new FileReader();
    
                reader.onload = (e) => resolve(e.target.result)
                reader.onerror = (e) => reject(reader.error)
        
                reader.readAsText(file)
            },
            oncancel: e => {
                resolve("")
            },
        })

        fileInput.click();
    });
}

const UploadJSON = async ({onSuccess, onErr}) => {
    let text
    try {
        text = await ReadFileInputAsText();
    } catch(err) {
        onErr({
            title: "Action Failed",
            message: "Failed to read file input",
            details: [err.message],
        })

        return
    }

    if (text.length === 0) {
        console.log("input file select canceled")

        return
    }

    try {
        const obj = JSON.parse(text)

        onSuccess(obj)
    } catch(err) {
        onErr({
            title: "Invalid Input",
            message: "Failed to parse input file as JSON",
            details: [err.message],
        })
    }
}

export {UploadJSON};