var Connection = new Client("ws://localhost:8080/ws");

var Messages = {
	list: [],
	init: function() {
		Connection.handle("publish", function(raw) {
			var msg = Object.assign({
				kind: "publish",
			}, raw)
			Messages.list.push(msg)
			m.redraw()
		})
		Connection.handle("rename", function(raw) {
			Messages.list.push(Object.assign({
				kind:"rename",
			}, raw))
			m.redraw()
		})
		Connection.handle("greet", function(raw) {
			Messages.list.push(Object.assign({
				kind:"greet",
			}, raw))
			m.redraw()
		})
		Connection.handle("goodbye", function(raw) {
			Messages.list.push(Object.assign({
				kind:"goodbye",
			}, raw))
			m.redraw()
		})
	},
	send: function(msg) {
		Connection.call("publish", Object.assign({}, msg, {
			time: "" + msg.time
		}))
	}
};

var User = {
	name: ""
}

var Chat = {
	lastScroll: 0,
	onupdate: function(vnode) {
		var scroll = vnode.dom.scrollHeight
		if (Chat.lastScroll == scroll) {
			return
		}
		Chat.lastScroll = scroll
		document.body.scrollTop = scroll
	},
	view: function() {
		return m("div.messages", [
			Messages.list.map(function(msg) {
				var d = new Date(msg.time);
				switch (msg.kind) {
					case "rename":
						return m("p.message", [
							m("span.message-time",   d.toLocaleTimeString()),
							m("span.message-prev",   msg.prev),
							m("span.message-invite", "~>"),
							m("span.message-last",   msg.name)
						])
					case "publish":
						return m("p.message", [
							m("span.message-time",   d.toLocaleTimeString()),
							m("span.message-author", msg.author),
							m("span.message-invite", ">"),
							m("span.message-text",   msg.text)
						])
					case "greet":
						return m("p.message", [
							m("span.message-time",   d.toLocaleTimeString()),
							m("span.message-author", msg.name),
							m("span.message-info",   "is here!")
						])
					case "goodbye":
						return m("p.message", [
							m("span.message-time",   d.toLocaleTimeString()),
							m("span.message-author", msg.name),
							m("span.message-info",   "gone =(")
						])
				}

			})
		])
	}
};

var App = {
	view: function(vnode) {
		var nav = function(route, caption) {
			var p = {
				role: "presentation",
			}
			if (m.route.get() == route) {
				p.className = "disabled"
			}
			return m("li", p, [ 
				m("a", { href: route, oncreate: m.route.link }, caption) 
			])
		}
		return [
			m("header.header", [
				m("div.container", [
					m(".row", [
						m(".col-xs-7", [
							m("ul.nav.nav-pills", [
								m("li.nav-header", {role: "presentation"}, [
									m("a#chat", { href: "/chat", oncreate: m.route.link }, "GoChat"),
								]),
								nav("/about", "about"),
							]),
						]),
						m(".col-xs-5.text-right", [
							m("form.user", {
								onsubmit: function(e) {
									e.preventDefault()
									if (User.name.length != 0) {
										Bootstrap.rename()
									}
								}	
							}, [
								m("span.glyphicon.glyphicon-user"),
								m("span.user-prefix", "@"),
								m("input.user-name", {
									value: User.name,
									oncreate: function(vnode) {
										vnode.dom.style.width = textWidth(vnode.dom.value)
									},
									onfocus: function(e) {
										var prev = this.value
										setTimeout(function() {
											e.target.value = prev
										}, 1)
									},
									oninput: function(e) {
										var el = e.target
										User.name = el.value
										el.style.width = textWidth(el.value)
									},
									onchange: function(e) { 
										if (User.name.length != 0) {
											Bootstrap.rename()
										}
									}
								}),
							])
						])
					])
				])
			]),
			m("div.container.content", [
				m("section", vnode.children)
			])
		]
	}
};

function textWidth(text) {
	var ret = 0;
	var temp = document.getElementById("temp");
	m.render(temp, m("div", {
		oncreate: function(vnode) {
			ret = vnode.dom.clientWidth;
		},
		onupdate: function(vnode) {
			ret = vnode.dom.clientWidth;
		},
		style: {
			"font-weight": "500",
			"font-size":   "14px",
			"position":    "absolute",
			"visibility":  "hidden",
			"height":      "auto",
			"width":       "auto",
			"white-space": "nowrap"
		},
	}, text))
	return ret + 5 + "px"
}

var Message = {
	text: "",
	reset: function() {
		var text = Message.text
		Message.text = ""
		return text
	}
}

var Compose = {
	oncreate: function() {
		setTimeout(function() {
			document.getElementById("compose").focus()
		}, 10)
	},
	view: function() {
		return m("footer.footer", [
			m("div.container", [
				m("form.form-horizontal.compose",
					{onsubmit: function(e) {
						e.preventDefault()
						var text = Message.reset()
						if (text.length == 0) {
							return
						}
						Messages.send({
							author: User.name,
							text:   text,
							time:   Date.now()
						})
					}},
				   	[
						m("div.form-group", [
							m("div.col-xs-12", [
								m("input.form-control.compose-input#compose", {
									type: "text",
									value: Message.text,
									placeholder: "Write a message...",
									autocomplete: "off",
									oninput: m.withAttr("value", function(value) { 
										Message.text = value
									})
								})
							])
						]),
					]
				)
			])
		])
	}
}

var About = {
	view: function() {
		return m("div", "hello, websocket!")
	}
};

var Bootstrap = {
	ready: false,
	oninit: function() {
		Messages.init()
		Connection.handle("hello", function(raw) {
			User.name = raw.name

			Bootstrap.ready = true
			if (Bootstrap.spinner) {
				Bootstrap.spinner.stop()
			}
			m.redraw();
		})
	},
	connect: function() {
		var self = this;
		this.ready = false

		m.route.set("/")

		return Connection.connect()
			.catch(function(err) {
				console.warn("connect error:", err);
				self.err = err;
			})
	},
	rename: function() {
		return Connection
			.call("rename", { name: User.name })
			.then(function() {
				console.log("rename ok" )
			})
			.catch(function(err) {
				console.warn("rename error:", err);
			});
	},
	oncreate: function() {
		if (this.ready) {
			return
		}
		if (this.err != null) {
			return;
		}

		setTimeout(function() {
			if (Bootstrap.ready) {
				return
			}
			var opts = {
				lines:   17,
				length:  12,
				width:   2,
				radius:  12,
				color:   '#268bd2',
				opacity: 0.1,
				speed:   1.5,
			}
			Bootstrap.spinner = new Spinner(opts).spin(document.body)
		}, 500)

		return Bootstrap.connect()
	},
	view: function(vnode) {
		if (this.ready) {
			m.route.set("/chat")
			return
		}
		if (this.err != null) {
			return m(".crash", [
				m(".crash-message", "Oh snap! Something went wrong! =(")
			])
		}
		return m("div", "loading...")
	}
}


m.route(document.body, "/", {
	"/": {
		render: function() {
			return m(Bootstrap)
		},
	},
	"/chat": {
		render: function() {
			if (!Bootstrap.ready) {
				m.route.set("/")
				return
			}
			return [ m(App, m(Chat)), m(Compose) ]
		}	
	},
    "/about": {
		render: function() {
			return m(App, m(About))
		}	
	}
});
