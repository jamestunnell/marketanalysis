import { Get } from '../backend.js'

function loadAllSettings() {
    return new Promise((resolve, reject) => {
        console.log('loading all settings')

        Get(`/settings`).then(resp => {
            if (resp.status != 200) {
                resp.json().then(appErr => {
                    console.log("failed to load all setting", appErr);

                    reject(appErr)
                })
        
                return;
            }

            resp.json().then(obj => {
                console.log("loaded all setting", obj.settings);

                resolve(obj.settings)
            })
        }).catch(err => {
            reject({
                title: "Action Failed",
                message: "failed to make get request",
                details: [err.message],
            })
        })
    })
}

export default loadAllSettings