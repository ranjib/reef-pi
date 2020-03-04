import React from 'react'
import { fetchDrivers, fetchDriverOptions, deleteDriver, createDriver, updateDriver } from 'redux/actions/drivers'
import { connect } from 'react-redux'
import Driver from './driver'
import New from './new'

class drivers extends React.Component {
  constructor (props) {
    super(props)
    this.list = this.list.bind(this)
    this.validate = this.validate.bind(this)
  }

  componentDidMount () {
    this.props.fetch()
    this.props.fetchDriverOptions()
  }

  validate (payload) {
    // This doesn't seem to belong in redux
    // since it isn't really part of app state.
    // It's here since I'm not sure where else it should belong.
    return fetch('api/drivers/validate', {
      method: 'POST',
      credentials: 'same-origin',
      body: JSON.stringify(payload)
    })
  }

  list () {
    const items = []
    this.props.drivers.sort((a, b) => { return parseInt(a.id) < parseInt(b.id) }).forEach((d, n) => {
      if (d.type === 'rpi') {
        items.push(
          <div className='row ' key={d.id}>
            <div className='col-4 col-md-6'>{d.name}</div>
            <div className='col-4 col-md-6'>{d.type}</div>
          </div>
        )
        return
      }
      items.push(
        <Driver
          key={d.id}
          driver={d}
          validate={this.validate}
          driverOptions={this.props.driverOptions}
          remove={this.props.delete}
          update={this.props.update}
        />
      )
    })
    return items
  }

  render () {
    return (
      <div className='container'>
        <div className='row mb-1'>
          <div className='col-12'>
            {this.list()}
          </div>
        </div>
        <New
          drivers={this.props.drivers}
          hook={this.props.create}
          driverOptions={this.props.driverOptions}
          validate={this.validate}
        />
      </div>
    )
  }
}

const mapStateToProps = state => {
  return {
    drivers: state.drivers,
    driverOptions: state.driverOptions
  }
}

const mapDispatchToProps = dispatch => {
  return {
    fetch: () => dispatch(fetchDrivers()),
    fetchDriverOptions: () => dispatch(fetchDriverOptions()),
    create: d => dispatch(createDriver(d)),
    delete: id => dispatch(deleteDriver(id)),
    update: (id, p) => dispatch(updateDriver(id, p))
  }
}

const Drivers = connect(
  mapStateToProps,
  mapDispatchToProps
)(drivers)
export default Drivers
