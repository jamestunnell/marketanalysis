class Output {
    constructor({name, measurements}) {
        this.name = name
        this.measurements = measurements
    }
}

function MakeOutputs({infos, configs}) {
    return infos.map(info => {
        const cfg = configs.find(cfg => cfg.name === info.name)
        
        return new Output({
            name: info.name,
            measurements: cfg ? cfg.measurements : [],
        })
    })
}

export {MakeOutputs}