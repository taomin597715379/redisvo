$(document).ready(function() {
	var host = $(location).attr('host');
	var addServer = new Vue({
		el: '#addServer',
		data:{
			name: '127.0.0.1',
			host: '127.0.0.1',
			port: 6379,
			auth: null,
			isIpValid: false,
			isPortValid: false,
			isIPAndPortValid: false,
		},
		methods: {
			inputIpOrPort: function() {
				this.isIpValid = false;
				this.isPortValid = false;
				this.isIPAndPortValid = false;
			},
			validateIp: function() {
				var exp1 = /^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/;
				var exp2 = /^\d+$/;
				var reg1 = this.host.match(exp1);
				var reg2 = this.port.match(exp2);
				if(reg1 == null && reg2 != null) {
					this.isIpValid = true;
				} else if(reg1 != null && reg2 == null) {
					this.isPortValid = true;
				} else if(reg1 == null && reg2 == null) {
					this.isIPAndPortValid = true;
				} else {
					var host = $(location).attr('host');
					$.ajax({
						'url':"http://"+host+"/addServer?name="+this.name+"&host="+this.host+"&port="+this.port+"&auth="+this.auth,
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
			removeElement: function (index) {
				this.my_tasks.splice(index, 1);
			},
			deleteTask: function(server_info, index) {
				var host = $(location).attr('host');
				$.ajax({
					'url':"http://"+host+"/getlist?status=remove&name="+server_info,
				  	'success':function(result, textStatus, request) {
				  		if(request.status == 200) {
				  			tasks.removeElement(index);
				  		}
				  	},
				  	'dataType':'text',
					'error':function() {
						console.log('error ajax ... ')
					}
				})
			}
		}
	});

	var version = new Vue({
		el: '#update-version',
		data: {
			is_update_version: false,
			now_version: null,
			update_version: null,
		},
		methods: {
			fetchData: function(is_need_version, now, update) {
				this.is_update_version = is_need_version;
				this.now_version = now;
				this.update_version = update;
				if(this.is_update_version == true) {
					$('#update-version').css("display", "block");
				} else {
					$('#update-version').css("display", "none");
				}
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
				modal.fetchServer(data.server_ext_infos);
			} else {
				modal.fetchModal(0);
				tasks.fetchData(data.server_ext_infos);
				version.fetchData(data.is_update_version, data.now_version, data.update_version);
			}
	  	},	
		'error':function() {
			console.log('error ajax ... ')
		}
	})
})



	

 