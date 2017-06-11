$(document).ready(function() {
	var host = $(location).attr('host');
	var submit = new Vue({
		el: '#submit',
		methods: {
		validateIp: function() {
			obj_ip  = $("#server-host").val();
			var exp = /^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/; 
			var reg = obj_ip.match(exp);
			if(reg == null) {
				alert("IP地址不合法！");
			} else {
				obj_name  = $("#server-name").val();
				obj_port  = $("#server-port").val();
				obj_auth  = $("#server-auth").val();
				var host = $(location).attr('host');
				$.ajax({
					'url':"http://"+host+"/addServer?name="+obj_name+"&host="+obj_ip+"&port="+obj_port+"&auth="+obj_auth,
				  	'success':function(result, textStatus, request) {
						location.reload();
				  	},
				  	'dataType':'text',
					'error':function() {
						console.log('error ajax ... ')
					}
				})
			}
		}
		}
	});
	var tasks = new Vue({
		el: '#servertbl',
		data: {
			my_tasks: [],
		},
		methods: {
			fetchData: function(server_info) {
				this.my_tasks = server_info;
			},
			deleteTask: function(server_info) {
				var host = $(location).attr('host');
				$.ajax({
					'url':"http://"+host+"/getlist?status=remove&name="+server_info,
				  	'success':function(result) {
				  		location.reload();
				  	},
				  	'dataType':'text',
					'error':function() {
						console.log('error ajax ... ')
					}
				})
			}
		}
	});

	var modal = new Vue({
	  el: '#admin-modal',
		data: {
			active: 0,
			registerAdmin: '',
			registerPassword: '',
			registerError: '',
			my_tasks: [],
		},
		methods: {
			fetchModal: function(enable) {
				this.active = enable;
			},
			fetchServer: function(enable) {
				this.my_tasks = enable;
			},
			submit: function(e) {
				e.preventDefault();
				if(this.registerAdmin != '' && this.registerPassword != '') {
					$.ajax({
						'url':"http://"+host+"/login?admin="+ this.registerAdmin +"&password="+this.registerPassword,
					  	'success': function(result, textStatus, request) {
					  		if(result == 'true' && request.status == 200) {
					  			Cookies.set('Auth', 'Restricted', { expires:1 });
					  			modal.registerError = '';
					  			$(".user-modal-container.active").removeClass("active");
					  			tasks.fetchData(modal.my_tasks);
					  		} else {
					  			modal.registerError = 'admin or password error';
					  		}
					  	},
					  	'dataType':'text',
						'error':function() {
							console.log('error ajax ... ')
						}
					})
				}
				if(this.registerAdmin == '') {
					modal.registerError = 'please input username';
				} else if (this.registerPassword == '') {
					modal.registerError = 'please input password';
				}
			}
		}
	});

	$.ajax({
		'url':"http://"+host+"/getlist?status=serverinfo",
	  	'dataType':'json',
	  	'success': function(data, textStatus, request) {
			if(request.getResponseHeader('Auth') == "Restricted") {
				modal.fetchModal(1);
				modal.fetchServer(data);
			} else {
				modal.fetchModal(0);
				tasks.fetchData(data);
			}
	  	},	
		'error':function() {
			console.log('error ajax ... ')
		}
	})
})



	

 