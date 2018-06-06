import React from 'react'
import Auth from '../auth.jsx'
import Capabilities from './capabilities.jsx'
import Display from './display.jsx'
import HealthNotify from './health_notify.jsx'
import {hideAlert} from '../utils/alert.js'
import {updateSettings, fetchCapabilities, fetchSettings} from '../redux/actions'
import {connect} from 'react-redux'

class settings extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      capabilities: props.capabilities,
      settings: {
        name: '',
        address: '',
        interface: ''
      },
      updated: false
    }
    this.updateCheckbox = this.updateCheckbox.bind(this)
    this.showCapabilities = this.showCapabilities.bind(this)
    this.updateCapabilities = this.updateCapabilities.bind(this)
    this.update = this.update.bind(this)
    this.showDisplay = this.showDisplay.bind(this)
    this.toRow = this.toRow.bind(this)
    this.updateHealthNotify = this.updateHealthNotify.bind(this)
    this.showHealthNotify = this.showHealthNotify.bind(this)
  }

  showHealthNotify () {
    if (this.state.settings.health_check === undefined) {
      return
    }
    if (this.state.settings.capabilities.health_check !== true) {
      return
    }
    return (
      <HealthNotify update={this.updateHealthNotify} state={this.state.settings.health_check} />
    )
  }

  updateHealthNotify (notify) {
    if (notify !== undefined) {
      var settings = this.state.settings
      settings.health_check = notify
      this.setState({settings: settings, updated: true})
    }
    return this.state
  }

  updateCheckbox (key) {
    return (function (ev) {
      var settings = this.state.settings
      settings[key] = ev.target.checked
      this.setState({
        settings: settings,
        updated: true
      })
    }.bind(this))
  }

  showDisplay () {
    if (!this.state.settings.display) {
      return
    }
    return (<div className='container'><Display /></div>)
  }

  showCapabilities () {
    return (
      <Capabilities capabilities={this.state.capabilities} update={this.updateCapabilities} />
    )
  }

  updateCapabilities (capabilities) {
    var settings = this.state.settings
    settings.capabilities = capabilities
    this.setState({
      settings: settings,
      updated: true
    })
  }

  update () {
    this.props.updateSettings(this.state.settings)
    this.setState({updated: false})
  }

  componentDidMount () {
    this.props.fetchCapabilities()
    this.props.fetchSettings()
  }

  toRow (label) {
    var fn = function (ev) {
      var settings = this.state.settings
      settings[label] = ev.target.value
      this.setState({
        settings: settings,
        updated: true
      })
    }.bind(this)
    return (
      <div className='input-group'>
        <label className='input-group-addon'> {label}</label>
        <input type='text' onChange={fn} value={this.state.settings[label]} id={'to-row-' + label} />
      </div>
    )
  }

  render () {
    var updateButtonClass = 'btn btn-outline-success col-sm-2'
    if (this.state.updated) {
      updateButtonClass = 'btn btn-outline-danger col-sm-2'
    }

    return (
      <div className='container'>
        <div className='row'>
          <div className='col-sm-6'>
            {this.toRow('name')}
            {this.toRow('interface')}
            {this.toRow('address')}
            <div className='input-group'>
              <label className='input-group-addon'>Notification</label>
              <input type='checkbox' id='updateNotification' onClick={this.updateCheckbox('notification')} defaultChecked={this.state.settings.notification} className='form-control' />
            </div>
            <div className='input-group'>
              <label className='input-group-addon'>Display</label>
              <input type='checkbox' id='updateDisplay' onClick={this.updateCheckbox('display')} defaultChecked={this.state.settings.display} className='form-control' />
              {this.showDisplay()}
            </div>
            <div className='input-group'>
              <label className='input-group-addon'>Use HTTPS</label>
              <input type='checkbox' id='use_https' onClick={this.updateCheckbox('https')} defaultChecked={this.state.settings.https} className='form-control' />
            </div>
            <div className='input-group'>
              <label className='input-group-addon'>Enable PCA9685</label>
              <input type='checkbox' id='enable_pca9685' onClick={this.updateCheckbox('pca9685')} defaultChecked={this.state.settings.pca9685} className='form-control' />
            </div>
          </div>
        </div>
        <div className='row'>
          <div className='col-sm-4' >
            <label> <b>Capabilities</b> </label>
            {this.showCapabilities()}
          </div>
        </div>
        <div className='row'>
          {this.showHealthNotify()}
        </div>
        <div className='row'>
          <input type='button' className={updateButtonClass} onClick={this.update} id='systemUpdateSettings' value='update' />
        </div>
        <div className='row'>
          <Auth />
        </div>
      </div>
    )
  }
}
const mapStateToProps = (state) => {
  return {
    capabilities: state.capabilities,
    settings: state.settings
  }
}

const mapDispatchToProps = (dispatch) => {
  return {
    fetchCapabilities: () => dispatch(fetchCapabilities()),
    fetchSettings: () => dispatch(fetchSettings()),
    updateSettings: (s) => dispatch(updateSettings(s)),
  }
}

const Settings = connect(mapStateToProps, mapDispatchToProps)(settings)
export default Settings
