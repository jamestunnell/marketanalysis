class IntParam {
    constructor({name, constraint, value}) {
        this.name = name
        this.constraint = constraint
        this.value = value
    }
}

class FloatParam {
    constructor({name, constraint, value}) {
        this.name = name
        this.constraint = constraint
        this.value = value
    }
}

function MakeParams({infos, values}) {
    const params = []
    
    infos.forEach(info => {
        switch (info.valueType) {
            case 'int':
                params.push(new IntParam({
                    name: info.name,
                    value: van.state(values[info.name] ?? info.defaultValue),
                    constraint: MakeConstraint(info.constraint),
                }))
        
                break
            case 'float64':
                params.push(new FloatParam({
                    name: info.name,
                    value: van.state(values[info.name] ?? info.defaultValue),
                    constraint: MakeConstraint(info.constraint),
                }))
        
                break
        }

        console.log(`unsupported param type ${info.valueType}`)
    })

    return params
}

export {MakeParams}