class Input {
    constructor({name, source}) {
        this.name = name
        this.source = source
    }
}

function MakeInputs({infos, configs}) {
    return infos.map(info => {
        const cfg = configs.find(cfg => cfg.name === info.name)
        
        return new Input({
            name: info.name,
            source: cfg ? cfg.source : "",
        })
    })
}

export {MakeInputs}