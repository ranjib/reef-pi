import React from 'react'
import PropTypes from 'prop-types'

export default class AutoProfile extends React.Component {
  constructor (props) {
    super(props)
    this.curry = this.curry.bind(this)
    this.sliderList = this.sliderList.bind(this)
    if (Array.isArray(props.config.values)){
      this.state = {
        values: props.config.values
      }
    }
    else{
      this.state = {
        values: Array(12).fill(0)
      }
    }    
  }

  curry (i) {
    return (ev) => {

      if (/^([0-9]{0,2}$)|(100)$/.test(ev.target.value)){
        var val = parseInt(ev.target.value)
        if (isNaN(val)) 
          val = ''
        
        var values = Object.assign(this.state.values)
        values[i] = val
        this.props.onChangeHandler({values: values})
        this.setState({values: values})
      }
    }
  }

  sliderList () {
    var values = Object.assign({}, this.state)
    values = values.values
   
    var rangeStyle = {
      WebkitAppearance: 'slider-vertical',
      writingMode: 'bt-lr',
      padding: '0 5px',
      width: '8px',
      height: '175px'
    }
    var list = []
    var labels = [
      '12 am',
      '2 am',
      '4 am',
      '6 am',
      '8 am',
      '10 am',
      '12 pm',
      '2 pm',
      '4 pm',
      '6 pm',
      '8 pm',
      '10 pm'
    ]

    for (var i = 0; i < 12; i++) {
      if (values[i] === undefined) {
        values[i] = 0
      }
      list.push(
        <div className='col-md-1 text-center' key={i + 1}>
          <div className="d-block d-md-none d-lg-block">
            <input type="text" name="value"
              className="form-control form-control-sm mb-1 d-block d-md-none d-lg-block"
              value={values[i]}
              onChange={this.curry(i)}
              disabled={this.props.readOnly} />
          </div>
          <div className="d-none d-md-inline d-lg-none">
            {values[i]}
          </div>
          <input
            className='d-none d-md-inline'
            type='range'
            style={rangeStyle}
            onChange={this.curry(i)}
            value={values[i]}
            id={'intensity-' + i}
            orient='vertical'
            disabled={this.props.readOnly}
          />
          <div>
            {labels[i]}
          </div>
        </div>
      )
    }
    return (list)
  }

  render () {
    return (
      <div className='container'>        
        <div className='row'>
          {this.sliderList()}
        </div>
      </div>
    )
  }
}

AutoProfile.propTypes = {
  config: PropTypes.object,
  onChangehandler: PropTypes.func,
  readOnly: PropTypes.bool
}
