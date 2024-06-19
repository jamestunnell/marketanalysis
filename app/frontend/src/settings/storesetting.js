import { PutJSON } from '../backend.js'

function storeSetting({name, value}) {
    return new Promise((resolve, reject) => {
        const object = {name, value}

        console.log("storing setting", object)

        PutJSON({route: `/settings/${name}`, object}).then(putResp => {
            if (putResp.status != 204) {
                putResp.json().then(appErr => {
                    console.log("failed to store setting", appErr);

                    reject(appErr)
                })
        
                return;
            }

            console.log("stored setting", object)

            // Avoid Fetch failed loading
            putResp.text().then(text => {
                resolve()
            })
        }).catch(err => {
            reject({
                title: "Action Failed",
                message: `failed to make put request`,
                details: [err.message],
            })
        })
    })
}

export default storeSetting