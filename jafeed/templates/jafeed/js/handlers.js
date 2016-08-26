if (event_class == '__jafeed__') {
	console.log("Jafeed event");
	$('#jafeed_widget').show();
	num_posts = $('#jafeed_num_posts').html();
	console.log("Num: "+num_posts)
	new_num = parseFloat(num_posts)+1;
	console.log("New num! "+new_num);
	$('#jafeed_num_posts').html(new_num);
	return false
}