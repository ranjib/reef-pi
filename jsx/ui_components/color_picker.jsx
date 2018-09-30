import React from 'react'
import { HuePicker } from 'react-color'
import PropTypes from 'prop-types'

class ColorPicker extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      expand: false
    }

    this.handleColorChange = this.handleColorChange.bind(this)
  }

  handleColorChange (e) {
    const event = {
      target: {
        name: this.props.name,
        value: e.hex
      }
    }
    this.props.onChangeHandler(event)
    this.setState({expand: false})
  }

  render () {
    if (this.state.expand === false) {
      return (
        <button
          disabled={this.props.readOnly}
          onClick={() => this.setState({expand: true})}
          style={{backgroundColor: this.props.color, color: this.props.color}}
          className='btn btn-secondary col-12'>
          Choose
        </button>
      )
    }
    return (
      <HuePicker name={this.props.name}
        className='mt-2'
        color={this.props.color}
        onChangeComplete={this.handleColorChange} />
    )
  }
}

ColorPicker.propTypes = {
  name: PropTypes.string.isRequired,
  color: PropTypes.string.isRequired,
  readOnly: PropTypes.bool,
  onChangeHandler: PropTypes.func.isRequired
}

export default ColorPicker
