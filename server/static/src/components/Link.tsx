import * as React from 'react';

interface LinkProps {
  active: boolean
  children: any
  onClick: Function
}

export default class Link extends React.Component<LinkProps, void> {
  render() {
    const {active, children, onClick} = this.props
    if (!active) {
      return (
        <span>{children}</span>
      )
    }
    return (
      <a
        href="#"
        onClick={(e) => {
          e.preventDefault()
          onClick(e)
        }}>
        {children}
      </a>
    )
  }
}
