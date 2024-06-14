package clients

const PORT_STAT_RES = `<html>

<head>
<title>Port Stistics</title>
<link rel="stylesheet" type="text/css" href="/style.css">
</head>

<body>
<center>

<fieldset>
<legend>Port Statistics</legend>
<form method="post" action="/port.cgi?page=stats">
<table>
	<tr>
		<th style="width:80px;">Port</th>
		<th style="width:80px;">State</th>
		<th style="width:80px;">Link Status</th>
		<th style="width:100px;">TxGoodPkt</th>
		<th style="width:100px;">TxBadPkt</th>
		<th style="width:100px;">RxGoodPkt</th>
		<th style="width:100px;">RxBadPkt</th>
	</tr>
	<tr>
		<td>Port 1</td>
		<td>Enable</td>
		<td>Link Up</td>
		<td>1</td>
		<td>11</td>
		<td>111</td>
		<td>1111</td>
	</tr>
	<tr>
		<td>Port 2</td>
		<td>Enable</td>
		<td>Link Up</td>
		<td>2</td>
		<td>22</td>
		<td>222</td>
		<td>2222</td>
	</tr>
	<tr>
		<td>Port 3</td>
		<td>Disable</td>
		<td>Link Up</td>
		<td>3</td>
		<td>33</td>
		<td>333</td>
		<td>3333</td>
	</tr>
	<tr>
		<td>Port 4</td>
		<td>Enable</td>
		<td>Link Up</td>
		<td>4</td>
		<td>44</td>
		<td>444</td>
		<td>4444</td>
	</tr>
	<tr>
		<td>Port 5</td>
		<td>Enable</td>
		<td>Link Up</td>
		<td>5</td>
		<td>55</td>
		<td>555</td>
		<td>5555</td>
	</tr>
	<tr>
		<td>Port 6</td>
		<td>Enable</td>
		<td>Link Up</td>
		<td>6</td>
		<td>66</td>
		<td>666</td>
		<td>6666</td>
	</tr>
	<tr>
		<td>Port 7</td>
		<td>Enable</td>
		<td>Link Down</td>
		<td>7</td>
		<td>77</td>
		<td>777</td>
		<td>7777</td>
	</tr>
	<tr>
		<td>Port 8</td>
		<td>Disable</td>
		<td>Link Down</td>
		<td>8</td>
		<td>88</td>
		<td>888</td>
		<td>8888</td>
	</tr>
	<tr>
		<td>Port 9</td>
		<td>Enable</td>
		<td>Link Up</td>
		<td>9</td>
		<td>99</td>
		<td>999</td>
		<td>9999</td>
	</tr>
</table>
<br style="line-height:50%">
<input type="submit" name="submit" value="   Clear   ">
<input type="hidden" name="cmd" value="stats">
</form>
</fieldset>
</center>
</body>
</html>`

const PORT_RES string = `<html>

<head>
<title>Port Setting</title>
<link rel="stylesheet" type="text/css" href="/style.css">
<script type="text/javascript">
</script>
</head>

<body>
<center>

<fieldset>
<legend>Port Setting</legend>
<form method="post" action="port.cgi">
	<table border="1">
	<tr>
		<th class=MidSize>Port</th>
		<th>State</th>
		<th>Speed/Duplex</th>
		<th>Flow Control</th>
	</tr>
	<tr>
		<td align="center">
			<select name="portid" multiple size="8" >
			<option value="0">Port 1
			<option value="1">Port 2
			<option value="2">Port 3
			<option value="3">Port 4
			<option value="4">Port 5
			<option value="5">Port 6
			<option value="6">Port 7
			<option value="7">Port 8
		</td>
		<td align="center">
			<select name="state" class=MidSize>
			<option value="1">Enable
			<option value="0">Disable
			</select>
		</td>
		<td align="center">
			<select name="speed_duplex" class=MidSize>
			<option value="0">Auto
			<option value="1">10M/Half
			<option value="2">10M/Full
			<option value="3">100M/Half
			<option value="4">100M/Full
			<option value="5">1000M/Full
			<option value="6">2500M/Full
			</select>
		</td>
		<td align="center">
			<select name="flow" class=MidSize>
			<option value="0">Off
			<option value="1">On
			</select>
		</td>
	</tr>
	</table>
<br style="line-height:50%">
<input type="submit" name="submit" value="   Apply   ">
<input type="hidden" name="cmd" value="port">
</form>
<hr>
<form method="post" action="port.cgi">
	<mockPortHandlertable border="1">
	<tr>
		<th class=MidSize>Port</th>
		<th>State</th>
		<th>Speed/Duplex</th>
		<th>Flow Control</th>
	</tr>
	<tr>
		<td align="center">
			<select name="portid" multiple size="1" >
			<option value="8">Port 9
		</td>
		<td align="center">
			<select name="state" class=MidSize>
			<option value="1">Enable
			<option value="0">Disable
			</select>
		</td>
		<td align="center">
			<select name="speed_duplex" class=MidSize>
			<option value="0">Auto
			<option value="4">100M/Full
			<option value="5">1000M/Full
			<option value="6">2500M/Full
			<option value="8">10G/Full
			</select>
		</td>
		<td align="center">
			<select name="flow" class=MidSize>
			<option value="0">Off
			<option value="1">On
			</select>
		</td>
	</tr>
	</table>
<br style="line-height:50%">
<input type="submit" name="submit" value="   Apply   ">
<input type="hidden" name="cmd" value="port">
</form>
<hr>
<br>
<table border="1">
  <tr>
    <th rowspan="2" width="90">Port</th>
    <th rowspan="2" width="90">State</th>
    <th colspan="2">Speed/Duplex</th>
    <th colspan="2">Flow Control</th>
  </tr>
  <tr>
    <th width="90">Config</th>
    <th width="90">Actual</th>
    <th width="90">Config</th>
    <th width="90">Actual</th>
  </tr>
  <tr>
    <td>Port 1</td>
    <td>Enable</td>
    <td>Auto</td>
    <td>1000Full</td>
    <td>Off</td>
    <td>Off</td>
  </tr>
  <tr>
    <td>Port 2</td>
    <td>Enable</td>
    <td>100 Full</td>
    <td>100 Full</td>
    <td>On</td>
    <td>On</td>
  </tr>
  <tr>
    <td>Port 3</td>
    <td>Enable</td>
    <td>100 Half</td>
    <td>100 Half</td>
    <td>On</td>
    <td>Off</td>
  </tr>
  <tr>
    <td>Port 4</td>
    <td>Enable</td>
    <td>10 Full</td>
    <td>10 Full</td>
    <td>Off</td>
    <td>Off</td>
  </tr>
  <tr>
    <td>Port 5</td>
    <td>Enable</td>
    <td>10 Half</td>
    <td>10 Half</td>
    <td>Off</td>
    <td>Off</td>
  </tr>
  <tr>
    <td>Port 6</td>
    <td>Enable</td>
    <td>2.5G Full</td>
    <td>2500Full</td>
    <td>On</td>
    <td>On</td>
  </tr>
  <tr>
    <td>Port 7</td>
    <td>Enable</td>
    <td>Auto</td>
    <td>Link Down</td>
    <td>On</td>
    <td>Off</td>
  </tr>
  <tr>
    <td>Port 8</td>
    <td>Disable</td>
    <td>Auto</td>
    <td>Link Down</td>
    <td>Off</td>
    <td>Off</td>
  </tr>
  <tr>
    <td>Port 9</td>
    <td>Enable</td>
    <td>10G Full</td>
    <td>10G Full</td>
    <td>On</td>
    <td>Off</td>
  </tr>
</table>
<br>
</fieldset>
</center>
</body>
</html>`

const INFO_RES = `<html>
<head>
<title>System Information</title>
<link rel="stylesheet" type="text/css" href="/style.css">
</head>

<body>
<center>

<fieldset>
<legend>System Info</legend>
<br>
<table>
  <tr>
    <th style="width:150px;">Device Model</th>
    <td style="width:250px;">WAMJHJ-8125MNG</td>
  </tr>
  <tr>
    <th>MAC Address</th>
    <td>1C:2A:12:34:56:78</td>
  </tr>
  <tr>
    <th>IP Address</th>
    <td>192.168.2.1</td>
  </tr>
  <tr>
    <th>Netmask</th>
    <td>255.255.255.0</td>
  </tr>
  <tr>
    <th>Gateway</th>
    <td>192.168.2.254</td>
  </tr>
  <tr>
    <th>Firmware Version</th>
    <td>V1.9</td>
  </tr>
  <tr>
    <th>Firmware Date</th>
    <td>Jan 03 2024</td>
  </tr>
  <tr>
    <th>Hardware Version</th>
    <td>V1.1</td>
  </tr>
</table>
<br>
</fieldset>
</center>
</body>
</html>`
