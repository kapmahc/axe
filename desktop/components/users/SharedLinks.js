import {ListGroup, ListGroupItem} from 'reactstrap'
import Icon from 'react-fontawesome'
import Link from 'next/link'
import {FormattedMessage} from 'react-intl'

export default() => (<ListGroup>
  {
    [
      {
        icon: 'sign-in',
        href: '/users/sign-in',
        label: 'nut.users.sign-in.title'
      }, {
        icon: 'user-plus',
        href: '/users/sign-up',
        label: 'nut.users.sign-up.title'
      }, {
        icon: 'key',
        href: '/users/forgot-password',
        label: 'nut.users.forgot-password.title'
      }, {
        icon: 'check-square',
        href: '/users/confirm',
        label: 'nut.users.confirm.title'
      }, {
        icon: 'unlock',
        href: '/users/unlock',
        label: 'nut.users.unlock.title'
      }, {
        icon: 'commenting',
        href: '/leave-words/new',
        label: 'nut.leave-words.new.title'
      }
    ].map((it, id) => (<ListGroupItem action={true} key={id}>
      <Icon name={it.icon}/>
      &nbsp;
      <Link href={it.href}>
        <FormattedMessage id={it.label}/>
      </Link>
    </ListGroupItem>))
  }
</ListGroup>)
