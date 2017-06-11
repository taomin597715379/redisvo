$(document).ready(function() {
	// register click-outside event
	Vue.directive('click-outside', {
	 	priority: 700,
	 	bind() {
		    this.event = function (event) {
		    	event.preventDefault();
		    	typeName.hideContextMenu();
		    	typeName.rollback();
		 	}
		    this.el.addEventListener('click', this.stopProp)
		    document.body.addEventListener('click',this.event)
	  	},  
	  	unbind() {
	  		this.el.removeEventListener('click', this.stopProp)
			document.body.removeEventListener('click',this.event)
	  	},
	  	stopProp(event) {event.stopPropagation() }
	});
	// register context-menu-item
	Vue.component('context-menu-item', {
		template: '#template-context-menu-item',
	});
	// register context-menu
	Vue.component('context-menu', {
		template: '#template-context-menu',
		data: function() {
	        return {
				type: null,
				name: null,
				index: null,
				upHere: false,
	       }
	    },
		methods: {
			chgwidht1: function() {
				this.upHere = true;
				$("#context-menu").css("width", 210);
			},
			chgwidht2: function() {
				this.upHere = false;
				$("#context-menu").css("width", 170);
			},
			modify: function() {
				$("#type-name tbody tr td:nth-child(2)").each(function(index){
				    if( $(this).hasClass('before-modify')) {
				    	var val = $(this).find('input').val();
				        typeName.keys[index].name = val;
				        $("#type-name tbody tr td:nth-child(2)").eq(index).removeClass('before-modify')
				    }
				})
				typeName.keys[this.index].name = '<input style="border-radius: 4px;"   \
												  name="name_modify" type="text" v-model="this.name" \
												  value="' +  this.name  +'">';
				$('#type-name tbody tr td:nth-child(2)').eq(this.index).addClass('before-modify');
				typeName.hideContextMenu();
			},
			remove: function() {
				typeName.deleteKey(this.type, this.name);
				typeName.hideContextMenu();
			},
			refresh: function() {
				typeName.getKeysByTypeName(this.type, this.name);
				typeName.hideContextMenu();
			},
			cancel: function() {
				typeName.hideContextMenu();
			}
		},
		events: {
        	'contextmenu-info': function(index, type, name) {
        		this.type = type;
				this.name = name;
				this.index = index;
        	}
        }
	});
    $('.navbar-right li a').click(function() {
        $('.navbar-right li.active').removeClass('active');
        $(this).parent('li').addClass('active');
    });
    $('.navbar-right li a')[0].click();
    $("#leftbar").css({
		position: "relative"
	}).resizable({
		    resizeHeight: false,
		    handles: 'e',
		    distance: 0,
		    minWidth: 250,
		    maxWidth: 450,
		    resize: function(e, ui) {
                $('#rightbar').css("width", 1170 - $('#leftbar').width() - $('#midbar').width());
            }
	});
	$("#midbar").css({
		position: "relative"
	}).resizable({
		    resizeHeight: false,
		    handles: 'e',
		    distance: 0,
		    minWidth: 250,
		    maxWidth: 450,
		    resize: function(e, ui) {
                $('#rightbar').css("width", 1170 - $('#leftbar').width() - $('#midbar').width());
            }
	});
	$('#dbdropdown > li > a').click(function(e){
	    $('#dbdropdown > li > a').removeClass('selected');
	    $(this).addClass('selected');
	    var db = $('#dbdropdown > li').find("a.selected").text().trim()
    	var server_info = getUrlParameter('server')
    	var host = $(location).attr('host');
		$.ajax({
			'url':"http://"+host+"/serverinfo?server=" + server_info + "&db=" + db,
		  	'success':function(result) {
	  			typeName.fetchData(result.typename);
	  			keysName.fetchData(result.keysnameswithtype);
	  			content.fetchData(result.content)
	  			if(result.keysnameswithtype.hasOwnProperty('keysname') && 
		  		   result.keysnameswithtype.keysname !=null &&  result.keysnameswithtype.keysname.length > 0) {
		  			content.fetchField(result.keysnameswithtype.keysname[0].name, 0);
		  		}
		  		if(result.keysnameswithtype.hasOwnProperty('selftypename') && 
		  		   result.keysnameswithtype.selftypename != null) {
					content.fetchTypeName(result.keysnameswithtype.selftypename);
		  		}
		  		if (result.hasOwnProperty('typename') && (result.typename.length > 0)) {
		  			if (result.typename[0].type == "string") {
		  				$("#midbar").hide();
		  				$("#rightbar").toggleClass('col-md-6 col-md-9');
		  			} else {
		  				// $("#rightbar").toggleClass('col-md-9 col-md-6');
		  				$("#midbar").show();
		  				keysName.fetchData(result.keysnameswithtype);
		  			}
		  		}
		  	},
		  	'dataType':'json',
			'error':function() {
				console.log('error ajax ... ')
			}
		});
	});
	$("#dbdropdown > li > a")[0].click();
	$(".faq-links").click(function(e){
		content.parseString($("#rightbar textarea").val());
	});
    $("#config").on("change",function(){
        $(".savetoconfigfile > button").removeClass('disabled', false);
        $(".apply > button").removeClass('disabled', false);
    }).on("keyup", function(){
    	$(".savetoconfigfile > button").removeClass('disabled', false);
        $(".apply > button").removeClass('disabled', false);
    });
	var height = $(window).height();
	$("#config").height( height-70);
	$( window ).bind("resize", function(){
		var height = $(window).height();
	    $("#config").height( height-70 );
	});
	$("#leftbar").height( height-70 );
	$(".fht-tbody").height( height-38-70 );
	$( window ).bind("resize", function(){
		var height = $(window).height();
	    $("#leftbar").height( height-70 );
	    $(".fht-tbody").height( height-38-70 );
	});
	$("#midbar").height( height-70 );
	$(".fht-tbody").height( height-38-70 );
	$( window ).bind("resize", function(){
		var height = $(window).height();
	    $("#midbar").height( height-70 );
	    $(".fht-tbody").height( height-38-70 );
	});
	$( window ).bind("resize", function(){
		var height = $(window).height();
	    $("#terminal").height( height-90 );
	});
	$("#rightbar").height( height-38-70 );
	$( window ).bind("resize", function(){
		var height = $(window).height();
	    $("#rightbar").height( height-38-70 );
	});
	$('textarea').height( height-38-70);
	$( window ).bind("resize", function(){
		var height = $(window).height();
	    $('textarea').height( height-38-70 );
	});

	$('#rightbar input').height( height-38-70);
	$( window ).bind("resize", function(){
		var height = $(window).height();
	    $('#rightbar input').height( height-38-70 );
	});

	$('#type-name').fixedHeaderTable({
		cloneHeadToFoot: true,
	});
	$('#keys-name').fixedHeaderTable({ 
		cloneHeadToFoot: true,
	});
	$('#rightbar').css("width", 1170 - $('#leftbar').width() - $('#midbar').width());
	$('#type-name').css("width", "")
	$('#keys-name').css("width", "")
	$('#type-name').css("margin-top", -39);
	$('#keys-name').css("margin-top", -39);
	$('#leftbar thead tr th:eq(1)').css('padding-right', 5);
	$('#midbar thead tr th:eq(2)').css('padding-right', 5);
	
	var logo = new Vue({
		el:'#navbar-nlogo',
		methods: {
			back: function() {
				window.location.href='./index.html';
			}
		}
	})
	var typeName = new Vue({
		el: '#leftbar',
	  	data: {
	  		viewMenu: false,
	  	 	keys: [],
	  	 	show_more_flag: 0,
	  	 	contextMenuWidth: null,
    		contextMenuHeight: null,
    		keyName: null,
    		index: null,
	  	},
	  	methods: {
	  		notify: function(index, type, name) {
	  			this.$broadcast('contextmenu-info', index, type, name)
	  		},
	  		rollback: function() {
	  			if(this.keyName != null) {
	  				this.keys[this.index].name = this.keyName;
	  			}
	  		},
	  		resizeInput: function() {
	  			$("input[name=name_modify]").keyup(function() {
					$(this).attr('size', $(this).val().length);
				}).each(function() {
					$(this).attr('size', $(this).val().length);
				});
	  		},
	  		modifySubmit: function(type) {
	  			var name = $("input[name=name_modify]").val();
	  			var server_info = getUrlParameter('server')
		    	var db = $('a.selected').children("span").text();
		    	var host = $(location).attr('host');
		    	if(name != "") {
		    		$.ajax({
						'url':"http://"+host+"/modify?server=" + server_info + "&db=" + db + "&style=" + type + "&oldname=" + this.keyName + "&newname=" + name,
					  	'success':function(result) {
					  		location.reload();
					  	},
					  	'dataType':'json',
						'error':function() {
							this.keys[this.index].name = this.keyName;
							console.log('error ajax ... ')
						}
					})
		    	} else {
		    		this.keys[this.index].name = this.keyName
		    	}
	  		},
	  		showContextMenu: function(index, key, event) {
	  			this.keyName = key.name;
	  			this.index = index;
	  			this.notify(index, key.type, key.name);
			    var menu = document.getElementById("context-menu");
			    if(!this.contextMenuWidth || !this.contextMenuHeight) {
			        menu.style.visibility = "hidden";
			        menu.style.display = "block";
			        this.contextMenuWidth = menu.offsetWidth;
			        this.contextMenuHeight = menu.offsetHeight;
			        menu.removeAttribute("style");
			    }
			    if((this.contextMenuWidth + event.layerX) >= window.innerWidth) {
			        menu.style.left = (10+event.layerX - this.contextMenuWidth) + "px";
			    } else {
			        menu.style.left = 10+event.layerX + "px";
			    }
			    if((this.contextMenuHeight + event.layerY) >= window.innerHeight) {
			        menu.style.top = (40 + event.layerY - this.contextMenuHeight) + "px";
			    } else {
			        menu.style.top = 40 + event.layerY + "px";
			    }  
			    menu.classList.add('active');
		    },
		    hideContextMenu: function(index, event) {
		    	document.getElementById("context-menu").classList.remove('active');
		    },
		   	fetchData: function(obj) {
		    	this.keys = obj
		    },
		    addkey: function() {
		    	var server_info = getUrlParameter('server')
				var dbno = $('a.selected').children("span").text();
				var type = $('#addtypename select').val();
				var key = $('#server-name').val();
				var host = $(location).attr('host');
				if(key !="" && host !="") {
					$.ajax({
						'url':"http://"+host+"/addkey?server=" + server_info + "&db=" + dbno + 
																"&style=" + type + "&name=" + key,
					  	'success':function(result) {
				  			searValue.search_value = key;
				  			searValue.searchTypeName();
				  			$('#addtypename').modal('hide');
					  	},
					  	'dataType':'json',
						'error':function() {
							console.log('error ajax ... ')
						}
					})
				}
		    },
		    deleteKey: function(type, name) {
		    	var server_info = getUrlParameter('server')
		    	var db = $('a.selected').children("span").text();
		    	var host = $(location).attr('host');
		    	$.ajax({
					'url':"http://"+host+"/deletekey?server=" + server_info + "&db=" + db + "&style=" + type + "&name=" + name,
				  	'success':function(result) {
				  		location.reload();
				  	},
				  	'dataType':'json',
					'error':function() {
						console.log('error ajax ... ')
					}
				})
		    },
		    getMoreTypeName: function() {
		    	var server_info = getUrlParameter('server')
				var db = $('a.selected').children("span").text();
				var search_input = $('#search_input').val();
				this.show_more_flag = this.show_more_flag + 1;
				if(search_input == '') {
					var host = $(location).attr('host');
					$.ajax({
						'url':"http://"+host+"/serverinfo?server=" + server_info + "&db=" + db + "&showmore=" + this.show_more_flag,
					  	'success':function(result) {
				  			typeName.fetchData(result.typename);
					  	},
					  	'dataType':'json',
						'error':function() {
							console.log('error ajax ... ')
						}
		    		})
				} else {
					this.show_more_flag = 0;
					searValue.searchTypeNameShowMore()
				}
		    },
		    getKeysByTypeName: function(type, name) {
		    	var activeIcon = $('#nav_show .active a').text().trim();
		    	if(activeIcon == "Content") {
			    	if(type == "string") {
			    		if($("#rightbar").hasClass('col-md-6')) {
							$("#rightbar").toggleClass('col-md-6 col-md-9');
			    		}
			    		$("#midbar").hide();
			    	} else {
			    		if(!$("#rightbar").hasClass('col-md-6')) {
							$("#rightbar").toggleClass('col-md-9 col-md-6');
			    		}
			    		$("#midbar").show();
			    	}
			    	var server_info = getUrlParameter('server')
			    	var db = $('a.selected').children("span").text();
			    	var host = $(location).attr('host');
			    	$.ajax({
						'url':"http://"+host+"/showfields?server=" + server_info + "&db=" + db + "&style=" + type + "&name=" + name,
					  	'success':function(result) {
					  		keysName.fetchData(result.keysnameswithtype);
					  		content.fetchData(result.content)
					  		if(result.keysnameswithtype.hasOwnProperty('keysname') && 
					  		   result.keysnameswithtype.keysname !=null &&  result.keysnameswithtype.keysname.length > 0) {
					  			content.fetchField(result.keysnameswithtype.keysname[0].name, 0);
					  		}
					  		if(result.keysnameswithtype.hasOwnProperty('selftypename') && 
					  		   result.keysnameswithtype.selftypename != null) {
								content.fetchTypeName(result.keysnameswithtype.selftypename);
					  		}
					  	},
					  	'dataType':'json',
						'error':function() {
							console.log('error ajax ... ')
						}
					})
		    	}
		    }
		}
	});	
	var keysName = new Vue({
		el: '#midbar',
	  	data: {
	  	 	keys: [],
	  	 	selftypename: null,
	  	 	show_more_flag: 0,
	  	 	keytitle: 'key',
	  	 	index: 'index',
	  	 	zset_index: null,
	  	 	zset_index_show: false,
	  	},
	  	methods: {
		    fetchData: function(obj) {
		    	if(obj != null && obj.hasOwnProperty('keysname') &&
		    		obj.keysname != null && obj.keysname.length == 0) {
		    		this.keys = [];
		    		this.keytitle = "key";
		    		this.index = "index";
		    	}
		    	if (obj.hasOwnProperty('keysname') && obj.keysname != null &&
		    		obj.hasOwnProperty('selftypename') && obj.selftypename != null) {
		    		this.keys = obj.keysname
			    	this.selftypename = obj.selftypename
			    	this.zset_index_show = false;
			    	if(obj.selftypename.type == "hash") {
						this.keytitle = "key";
						this.index = "index";
					}
					if(obj.selftypename.type == "list") {
						this.keytitle = "item";
						this.index = "index";
					}
					if(obj.selftypename.type == "zset") {
						this.keytitle = "zmember";
						this.zset_index = "index";
						this.zset_index_show = true;
						this.index = "score";
					}
					if(obj.selftypename.type == "set") {
						this.keytitle = "member";
						this.index = "index";
					}
		    	}
		    },
		    addKeyByTypeAndName:function() {
		    	var server_info = getUrlParameter('server')
		    	var db = $('a.selected').children("span").text();
		    	var insert_flag = '';
		    	var add_value = '';
		    	var host = $(location).attr('host');
		    	switch (this.keytitle) {
		    		case "key":
						add_value = $('#key input').val();
						break;
		    		case "item":
		    			add_value = $('#item input').val();
		    			insert_flag = $('#item select').val();
		    			add_value = insert_flag + "_" + add_value;
		    			break;
		    		case "member":
		    			add_value = $('#member input').val();
		    			break;
		    		case "zmember":
		    			add_value = $("#zmember input[name='value']").val();
		    			insert_flag = $("#zmember input[name='score']").val();
		    			add_value = insert_flag + "_" + add_value;
		    			break;
		    		default:
		    			return;
		    	}
		    	if(add_value != '' && host != '') {
		    		$.ajax({
						'url':"http://"+host+"/addkey?server=" + server_info + "&db=" + db + "&style=" + this.selftypename.type + "&name=" + this.selftypename.name+ "&field=" + add_value,
					  	'success':function(result) {
					  		searValue.search_value = keysName.selftypename.name;
				  			searValue.searchTypeName();
				  			if(keysName.keytitle == "key") {
				  				$("#key").modal('hide');
				  			}
				  			if(keysName.keytitle == "item") {
				  				$("#item").modal('hide');
				  			}
				  			if(keysName.keytitle == "member") {
				  				$("#member").modal('hide');
				  			}
				  			if(keysName.keytitle == "zmember") {
				  				$("#zmember").modal('hide');
				  			}
					  	},
					  	'dataType':'json',
						'error':function() {
							console.log('error ajax ... ')
						}
					})
		    	}
		    },
		    getContentByKeys:function(key_name, index) {
		    	var server_info = getUrlParameter('server')
		    	var db = $('a.selected').children("span").text();
		    	var host = $(location).attr('host');
		    	content.fetchTypeName(this.selftypename);
		    	content.fetchField(key_name, index);
		    	$.ajax({
					'url':"http://"+host+"/showfields?server=" + server_info + "&db=" + db + "&style=" + this.selftypename.type + "&name=" + this.selftypename.name + "&key_name=" + key_name,
				  	'success':function(result) {
				  		content.fetchData(result.content)
				  	},
				  	'dataType':'json',
					'error':function() {
						console.log('error ajax ... ')
					}
				})
		    },
		   	getMoreFieldName: function() {
		    	var server_info = getUrlParameter('server')
				var db = $('a.selected').children("span").text();
				var host = $(location).attr('host');
				this.show_more_flag = this.show_more_flag + 1;
				$.ajax({
					'url':"http://"+host+"/showfields?server=" + server_info + "&db=" + db + "&style=" + this.selftypename.type + "&name=" + this.selftypename.name + "&showmore=" + this.show_more_flag,
				  	'success':function(result) {
				  		keysName.fetchData(result.keysnameswithtype);
				  	},
				  	'dataType':'json',
					'error':function() {
						console.log('error ajax ... ')
					}
	    		})
			}
		}
	});
	var content = new Vue({
		el: '#rightbar',
	  	data: {
	  		typename: {},
	  		field: "",
	  		index: 0,
	  		message:"",
	  	},
	  	mounted() {
	  		hljs.initHighlightingOnLoad();
	  	},
	  	methods: {
		    fetchData: function(obj) {
		    	$("#rightbar > a:nth-child(1)").hide();
		    	if(obj == "[]" || obj == "" ) {
		    		this.message = "";
		    	} else {
		    		try {
						var jsonObj = jQuery.parseJSON(obj);
						if(typeof jsonObj =='object') {
							var jsonPretty = JSON.stringify(jsonObj, undefined, 4);
			    			this.message = jsonPretty;
						} else {
							this.message = obj;
						}
		    		}
		    		catch(err) {
		    			this.message = obj;
		    		}
		    	}
		    },
		    fetchTypeName: function(typename) {
		    	this.typename = typename;
		    },
		    fetchField: function(Field, index) {
		    	this.field = Field;
		    	this.index = index;
		    },
		    parseString: function(obj) {
		    	if(obj == "[]" || obj == "" ) {
		    		this.message = "";
		    	} else {
					var jsonflag = $('.faq-links').find('i').hasClass('fa fa-plus-square-o');
					if(jsonflag) {
						$(".faq-links").find('i').removeClass('fa fa-plus-square-o');
						$('.faq-links').find('i').addClass('fa fa-minus-square-o');
						try {
							var jsonObj = jQuery.parseJSON(obj);
							if(typeof jsonObj =='object') {
								var jsonPretty = JSON.stringify(jsonObj, null, 0);
				    			this.message = jsonPretty;
							} else {
								this.message = obj;
							}
			    		}
			    		catch(err) {
			    			this.message = obj;
			    		}
					} else {
						$('.faq-links').find('i').removeClass('fa fa-minus-square-o');
						$('.faq-links').find('i').addClass('fa fa-plus-square-o');
						try {
							var jsonObj = jQuery.parseJSON(obj);
							if(typeof jsonObj =='object') {
								var jsonPretty = JSON.stringify(jsonObj, undefined, 4);
				    			this.message = jsonPretty;
							} else {
								this.message = obj;
							}
			    		}
			    		catch(err) {
			    			this.message = obj;
			    		}
					}
				}
		    },
		    changeContent: function() {
			    $("#rightbar > a:nth-child(1)").show();
		    },
		    saveContent:function() {
		    	var server_info = getUrlParameter('server')
				var db = $('a.selected').children("span").text();
				var host = $(location).attr('host');
		    	$.ajax({
					'url':"http://"+host+"/changeContent?server=" + server_info + "&db=" + db + "&style=" + this.typename.type + "&name=" + this.typename.name +
						"&index=" + this.index + "&field=" + this.field + '&content=' + this.message,
				  	'success':function(result) {
				  		$("#rightbar > a:nth-child(1)").hide();
				  	},
				  	'dataType':'json',
					'error':function() {
						console.log('error ajax ... ')
					}
	    		})
		    }
		}
	});
	var searValue = new Vue({
		el: '#search_input',
	  	data: {
	  	 	search_value: '',
	  	 	show_more_flag: 0,
	  	},
		methods: {
			searchTypeName: function() {
				$(document).keypress(function(e) {
				    if(e.which == 13) {
				    	e.preventDefault();
      					return ;
				    }
				});
				var db = $('a.selected').children("span").text();
				var server_info = getUrlParameter('server')
				var host = $(location).attr('host');
				$.ajax({
					'url':"http://"+host+"/search?server=" + server_info + "&db=" + db + "&search=" + this.search_value,
				  	'success':function(result) {
				  		typeName.fetchData(result.typename);
				  		var activeIcon = $('#nav_show .active a').text().trim();
		    			if(activeIcon == "Content") {
							keysName.fetchData(result.keysnameswithtype);
				  			content.fetchData(result.content)
				  			if(result.keysnameswithtype.hasOwnProperty('keysname') && 
					  		   result.keysnameswithtype.keysname !=null &&  result.keysnameswithtype.keysname.length > 0) {
					  			content.fetchField(result.keysnameswithtype.keysname[0].name, 0);
					  		}
					  		if(result.keysnameswithtype.hasOwnProperty('selftypename') && 
					  		   result.keysnameswithtype.selftypename != null) {
								content.fetchTypeName(result.keysnameswithtype.selftypename);
					  		}
				  			if (result.hasOwnProperty('typename') && (result.typename.length > 0)) {
					  			if (result.typename[0].type == "string") {
					  				$(".col-md-2").hide();
					  			} else {
					  				$(".col-md-2").show();
					  				keysName.fetchData(result.keysnameswithtype);
					  			}
					  		}
		    			}
				  	},
				  	'dataType':'json',
					'error':function() {
						console.log('error ajax ... ')
					}
				})
			},
			searchTypeNameShowMore: function() {
				$(document).keypress(function(e) {
				    if(e.which == 13) {
				    	e.preventDefault();
      					return ;
				    }
				});
				var db = $('a.selected').children("span").text();
				var server_info = getUrlParameter('server')
				var host = $(location).attr('host');
				this.show_more_flag = this.show_more_flag + 1;
				$.ajax({
					'url':"http://"+host+"/search?server=" + server_info + "&db=" + db + "&search=" + this.search_value + "&showmore=" + this.show_more_flag,
				  	'success':function(result) {
				  		typeName.fetchData(result.typename);
				  		var activeIcon = $('#nav_show .active a').text().trim();
		    			if(activeIcon == "Content") {
		    				keysName.fetchData(result.keysnameswithtype);
					  		content.fetchData(result.content);
					  		if(result.keysnameswithtype.hasOwnProperty('keysname') && 
					  		   result.keysnameswithtype.keysname !=null &&  result.keysnameswithtype.keysname.length > 0) {
					  			content.fetchField(result.keysnameswithtype.keysname[0].name, 0);
					  		}
					  		if(result.keysnameswithtype.hasOwnProperty('selftypename') && 
					  		   result.keysnameswithtype.selftypename != null) {
								content.fetchTypeName(result.keysnameswithtype.selftypename);
					  		}
					  		if (result.hasOwnProperty('typename') && (result.typename.length > 0)) {
					  			if (result.typename[0].type == "string") {
					  				$(".col-md-2").hide();
					  			} else {
					  				$(".col-md-2").show();
					  				keysName.fetchData(result.keysnameswithtype);
					  			}
					  		}
		    			}
				  	},
				  	'dataType':'json',
					'error':function() {
						console.log('error ajax ... ')
					}
				})
			}
		}
	});
	var conf = new Vue({
		el: '#config',
	  	data: {
	  		configValue:{},
	  	},
	  	validators:{
            numberical: function (val) {
            	if(val.length == 0) {
            		val = 0;
            	}
                return /^[0-9]\d*$/.test(val);
            }
	  	},
	  	methods: {
		    fetchData: function(obj) {
		    	this.configValue = obj;
		    },
		    saveChange: function() {
		    	var server_info = getUrlParameter('server')
		    	var host = $(location).attr('host');
		    	var self = this;
                self.$validate(true, function() {
                    if(self.$validatorMethod.invalid){
                        alert("Parameter Illegal");
                        location.reload();
                    } else {
                    	$.ajax({
							'type':'POST',
							'url':"http://"+host+"/config?server=" + server_info + "&operation=saveconfig",
							'data': JSON.stringify(self.configValue),
						  	'success':function(result) {
						  		location.reload();
						  	},
						  	'dataType':'json',
							'error':function() {
								location.reload();
								console.log('error ajax ... ')
							}
	    				})
                    }
                });
		    }
		}
	});
	var show = new Vue({
		el: '#nav_show',
	  	data: {
	  	},
	  	methods: {
		    terminalshow: function() {
				$("#midbar").css('display','none')
				$("#rightbar").css('display','none')
				$("#config").css('display','none')
				var server_info = getUrlParameter('server')
				var connection = null;
				var host = $(location).attr('host');
				$('#terminal').terminal(function(command, term) {
				    term.pause();
					// websocket
					window.WebSocket = window.WebSocket || window.MozWebSocket;
					connection = new WebSocket("ws://"+host+"/monitor");
			        if (command !== '') {
			        	if( (command == "monitor") || (command.startsWith("monitor") && command.indexOf("|") > 0) ) {
			        		if(command == "monitor") {
				        		term.set_prompt('');
				        		connection.onopen = function() {
				        			connection.send(server_info);
				        		};
				        	} else {
			        			var substr = "";
			        			term.set_prompt('');
								connection.onopen = function() {
				        			if (command.indexOf("findstr") > 0) {
				        				substr = command.substr(command.indexOf("findstr"))
				        			} else if(command.indexOf("grep") > 0) {
				        				substr = command.substr(command.indexOf("grep"))
				        			} else {
				        				term.echo("Pipeline commands should contain <findstr or grep>").resume();
				        				connection.close();
				        				term.echo("Bye... ").resume();
				        				term.set_prompt('iRedis> ');
				        			}
			        				connection.send(server_info + '_' + substr);
			        			};
			        		}
			        		connection.onmessage = function(response) {
						        term.echo(response.data).resume();
						    };
							connection.onclose = function() {
								term.echo("Bye... ").resume();
								term.set_prompt('iRedis> ');
							};
							connection.onerror = function() {
								term.echo("Bye... ").resume();
								term.set_prompt('iRedis> ');
							};
			        	} else {
			        		try {
				                $.ajax({
									'url':"http://"+host+"/terminal?command=" + command + "&server=" + server_info,
								  	'success':function(result) {
								  		term.echo(result).resume();
								  	},
								  	'dataType':'text',
									'error':function() {
										term.echo(new String('error ajax ... ')).resume();
									}
								})
				            } catch(e) {
				            	term.echo(new String(e)).resume();
				            }
			        	}
			        } else {
			        	term.echo('').resume();
			        }
			    }, {
			        name: 'js_demo',
			        height: 400,
			        prompt: 'iRedis> ',
			        enabled: true,
			        keypress: function(e,term) {
		        		if((e.keyCode == 26 || e.keyCode == 3) && e.ctrlKey) {
		        			if((connection !== null ) && connection.readyState !== connection.CLOSED &&
		        				connection.readyState !== connection.CLOSING){
				    			connection.send("Ctrl+C");
				    			connection.close();
				    			e.key = "";
				    			term.set_prompt('iRedis> ');
				    		}
		        		}
			        }
			    });
			    var height = $(window).height()-38-70;
				$("#terminal").css("height", height + 'px')
				$("#terminal").css("display", "inline")
		    },
		    contentshow: function() {
		    	$("#midbar").show()
				$("#rightbar").show()
				$("#terminal").css('display','none')
				$("#config").css('display','none')
		    },
		    configshow: function() {
		    	$("#midbar").css('display','none')
				$("#rightbar").css('display','none')
				$("#terminal").css('display','none')
				$("#config").show()
				var server_info = getUrlParameter('server')
				var host = $(location).attr('host');
				$.ajax({
					'url':"http://"+host+"/config?server=" + server_info + "&operation=showconfig",
				  	'success':function(result) {
						conf.fetchData(result)
				  	},
				  	'dataType':'json',
					'error':function() {
						console.log('error ajax ... ')
					}
				})
		    }
		}
	});
	show.contentshow();
});

var getUrlParameter = function getUrlParameter(sParam) {
    var sPageURL = decodeURIComponent(window.location.search.substring(1)),
        sURLVariables = sPageURL.split('&'),
        sParameterName,
        i;

    for (i = 0; i < sURLVariables.length; i++) {
        sParameterName = sURLVariables[i].split('=');

        if (sParameterName[0] === sParam) {
            return sParameterName[1] === undefined ? true : sParameterName[1];
        }
    }
};

