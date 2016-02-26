import * as React from 'react';

interface Props {
  indent: number;
  className: string;
  active: boolean
  title: string
  children: any
  onClick: Function
}

export default class Link extends React.Component<Props, void> {
  render() {
    const {indent, className, active, title, children, onClick} = this.props
    if (!active) {
      return (
        <span>{children}</span>
      )
    }
    return (
      <a
        className={className}
        style={{
          paddingLeft: (1.5 * indent) + 'em'
        }}
        href="#"
        title={title}
        onClick={(e) => {
          e.preventDefault()
          onClick(e)
        }}>
        {children}
      </a>
    )
  }
}
