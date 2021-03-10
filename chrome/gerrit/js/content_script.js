var contact='Please contact church.zhong@hmdglobal.com';
var url=window.location.href;
function bugOn() {console.log(contact);}
if(typeof url == "undefined" || null == url) {
	bugOn();
}
else
{
	//console.log('url=' + url);
	//hmdgerritserver.hmdglobal.com
	var reg = /[https|http]\:\/\/([^\/]*)\/([^\/]*)\/([^\/]*)/;
	var match = url.match(reg);
	//console.log('match=' + match);

	if(typeof match == "undefined" || null == match)
	{
		bugOn();
	}
	else
	{
		var host=match[1];
		var directory=match[2];
		var page=match[3];
		console.log('host=' + host);
		console.log('directory=' + directory);
		console.log('page=' + page);
	}
}
