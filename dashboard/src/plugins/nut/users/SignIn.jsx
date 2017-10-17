import React, {Component} from 'react'
import {FormattedMessage} from 'react-intl'

class Widget extends Component {
  render() {
    return (
      <div>
        <FormattedMessage id="nut.users.sign-in.title"/>
      </div>
    );
  }
}

export default Widget;
