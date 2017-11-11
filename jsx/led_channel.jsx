import React from 'react'
import LightSlider from './light_slider.jsx'

export default class LEDChannel extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      channel: this.props.ch,
    }
    this.sliderList = this.sliderList.bind(this)
    this.curry = this.curry.bind(this)
    this.updateAuto = this.updateAuto.bind(this)
    this.updateReverse = this.updateReverse.bind(this)
    this.updateFixedValue = this.updateFixedValue.bind(this)
    this.getFixedValue = this.getFixedValue.bind(this)
    this.update = this.update.bind(this)
  }

  update(k,v) {
    var ch = this.state.channel
    ch[k] = v
    this.setState({
      channel: ch
    })
    this.props.updateChannel(ch)
  }

  updateAuto (ev) {
    this.update('auto', ev.target.checked)
  }

  updateReverse (ev) {
    this.update('reverse', ev.target.checked)
  }

  getFixedValue () {
    return this.state.channel.fixed
  }

  updateFixedValue (v) {
    this.update('fixed', v)
  }

  curry (i) {
    return (
      function (ev) {
        var values  = this.state.channel.values
        values[i] = parseInt(ev.target.value)
        this.update('values', values)
      }.bind(this)
    )
  }

  sliderList () {
    var values = this.state.channel.values
    var rangeStyle = {
      WebkitAppearance: 'slider-vertical'
    }
    var list = []
    var labels = ['12 am', '2 am', '4 am', '6 am', '8 am', '10 am', '12 pm', '2 pm', '4 pm', '6 pm', '8 pm', '10 pm']

    for (var i = 0; i < 12; i++) {
      list.push(
        <div className='col-sm-1 text-center' key={i + 1}>
          <div className='row'>{values[i]}</div>
          <div className='row'>
            <input className='col-xs-1' type='range' style={rangeStyle} onChange={this.curry(i)} value={values[i]} id={'intensity-' + i} />
          </div>
          <div className='row'>
            <label>{labels[i]}</label>
          </div>
        </div>
      )
    }
    return (list)
  }

  render () {
    var showOnDemandSlider = {
      display: this.state.channel.auto ? 'none' : 'block'
    }
    var show24HourSliders = {
      display: this.state.channel.auto ? 'block' : 'none'
    }
    return (
      <div className='container'>
        <div className='row'>
          <div className='col-sm-2'>{this.props.name}</div>
        </div>
        <div className='row'>
          Auto<input type='checkbox' onClick={this.updateAuto} defaultChecked={this.state.channel.auto} id={this.props.name+'-auto'}/>
          Reverse<input type='checkbox' onClick={this.updateReverse} defaultChecked={this.state.channel.reverse} />
        </div>
        <div className='row' style={show24HourSliders}>
          {this.sliderList()}
        </div>
        <div className='row' style={showOnDemandSlider}>
          <LightSlider pin={this.props.pin} name={this.props.name} onChange={this.updateFixedValue} getValue={this.getFixedValue} style={showOnDemandSlider} />
        </div>
      </div>
    )
  }
}
