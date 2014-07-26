(function(
	Engine,
	Point,
	Polygon,
	Vector
){

Engine.Shape = function(x, y, width, height, points, polygons){
	var i, ref, point, poly;

	this.pos  = new Vector(x, y);
	this.size = new Vector(width, height);
	this.sizeRef = this.size.clone();

	ref = {};
	this.points = [];
	this.polygons = [];

	for (i = 0; i < points.length; i++) {
		point = new Point(
			points[i].id,
			points[i].x,
			points[i].y,
			this.size
		);
		ref[point.id] = point;
		this.points.push(point);
	}

	for (i = 0; i < polygons.length; i++) {
		poly = polygons[i];
		this.polygons.push(new Polygon(
			ref[poly.points[0]],
			ref[poly.points[1]],
			ref[poly.points[2]],
			poly.color
		));
	}
};

Engine.Shape.prototype = {

	breathing: false,

	breath: 0,
	breathLength: 1,
	breatheIn: false,

	startBreathing: function(){
		var p;

		this.breathing = true;
		this.breath = this.breathLength;

		for (p = 0; p < this.points.length; p++) {
			this.points[p].updateBreathingPhysics();
		}
	},

	breathe: function(tick){
		var p, scale, newSize;

		this.breath += tick;

		if (this.breath < this.breathLength) {
			return;
		}

		if (this.breatheIn) {
			scale = 1;
		} else {
			scale = 1.05;
		}

		this.breatheIn = !this.breatheIn;

		newSize = Vector.mult(this.sizeRef, scale);

		for (p = 0; p < this.points.length; p++) {
			this.points[p].updateTarget(newSize);
		}

		this.breath = 0;
	},

	update: function(engine){
		var p;

		if (this.breathing === true) {
			this.breathe(engine.tick);
		}

		for (p = 0; p < this.points.length; p++)  {
			this.points[p].update(engine);
		}

		for (p = 0; p < this.polygons.length; p++) {
			this.polygons[p].update(engine);
		}

		return this;
	},

	draw: function(ctx, scale, engine){
		var p, poly;

		ctx.translate(
			this.pos.x * scale >> 0,
			this.pos.y * scale >> 0
		);
		for (p = 0; p < this.polygons.length; p++) {
			poly = this.polygons[p];
			ctx.beginPath();
			ctx.moveTo(
				poly.a.pos.x * scale >> 0,
				poly.a.pos.y * scale >> 0
			);
			ctx.lineTo(
				poly.b.pos.x * scale >> 0,
				poly.b.pos.y * scale >> 0
			);
			ctx.lineTo(
				poly.c.pos.x * scale >> 0,
				poly.c.pos.y * scale >> 0
			);
			ctx.closePath();
			ctx.fillStyle   = poly.fillStyle;
			ctx.fill();
			ctx.lineWidth = 0.25 * scale;
			ctx.strokeStyle = poly.strokeStyle;
			ctx.stroke();
		}
		ctx.setTransform(1, 0, 0, 1, 0, 0);
		ctx.translate(
			engine.width  / 2 * engine.scale >> 0,
			engine.height / 2 * engine.scale >> 0
		);
		return this;
	}

};

})(
	window.Engine,
	window.Engine.Point,
	window.Engine.Polygon,
	window.Vector
);
