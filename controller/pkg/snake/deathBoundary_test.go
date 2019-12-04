package snake

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeathBoundary_Remove(t *testing.T) {
	type fields struct {
		deathBoundaries deathBoundary
	}
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name         string
		fields       fields
		args         []args
		expectations func(t *testing.T, s *Snake)
	}{
		{
			name: "removes a boundary point, but there's still y left",
			fields: fields{
				deathBoundaries: deathBoundary{
					5: {10: wallExists{}, 11: wallExists{}},
				},
			},
			args: []args{{5, 10}},
			expectations: func(t *testing.T, s *Snake) {
				assert.Equal(t, deathBoundary{
					5: {11: wallExists{}},
				}, s.deathBoundaries)
			},
		},
		{
			name: "removes a boundary x point, but no more x left",
			fields: fields{
				deathBoundaries: deathBoundary{
					5:  {10: wallExists{}},
					99: {1: wallExists{}},
				},
			},
			args: []args{{5, 10}},
			expectations: func(t *testing.T, s *Snake) {
				assert.Equal(t, deathBoundary{
					5:  {},
					99: {1: wallExists{}},
				}, s.deathBoundaries)
			},
		},
		{
			name: "tries to remove something that isn't found",
			fields: fields{
				deathBoundaries: deathBoundary{
					5:  {},
					99: {1: wallExists{}},
				},
			},
			args: []args{{100, 100}},
			expectations: func(t *testing.T, s *Snake) {
				assert.Equal(t, deathBoundary{
					5:  {},
					99: {1: wallExists{}},
				}, s.deathBoundaries)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Snake{
				deathBoundaries: tt.fields.deathBoundaries,
			}

			for _, arg := range tt.args {
				s.deathBoundaries.Remove(arg.x, arg.y)
			}

			tt.expectations(t, s)
		})
	}
}

func TestDeathBoundary_Add(t *testing.T) {
	type fields struct {
		deathBoundaries deathBoundary
	}
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name         string
		fields       fields
		args         []args
		expectations func(t *testing.T, s *Snake)
	}{
		{
			name: "adds a boundary when empty",
			fields: fields{
				deathBoundaries: deathBoundary{},
			},
			args: []args{{5, 10}},
			expectations: func(t *testing.T, s *Snake) {
				assert.Equal(t, deathBoundary{
					5: {10: wallExists{}},
				}, s.deathBoundaries)
			},
		},
		{
			name: "adds multiple",
			fields: fields{
				deathBoundaries: deathBoundary{},
			},
			args: []args{
				{5, 10},
				{4, 10},
				{3, 10},
			},
			expectations: func(t *testing.T, s *Snake) {
				assert.Equal(t, deathBoundary{
					5: {10: wallExists{}},
					4: {10: wallExists{}},
					3: {10: wallExists{}},
				}, s.deathBoundaries)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Snake{
				deathBoundaries: tt.fields.deathBoundaries,
			}

			for _, arg := range tt.args {
				s.deathBoundaries.Add(arg.x, arg.y)
			}

			tt.expectations(t, s)
		})
	}
}

func TestDeathBoundary_IsBoundary(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name string
		b    deathBoundary
		args args
		want bool
	}{
		{
			name: "can find a boundary easily",
			b: deathBoundary{
				1: {1: wallExists{}},
				2: {2: wallExists{}},
			},
			args: args{1, 1},
			want: true,
		},
		{
			name: "if it x doesn't exist, it's not a boundary",
			b: deathBoundary{
				1: {1: wallExists{}},
				2: {2: wallExists{}},
			},
			args: args{99, 99},
			want: false,
		},
		{
			name: "if it y doesn't exist, it's not a boundary",
			b: deathBoundary{
				1: {1: wallExists{}},
				2: {2: wallExists{}},
			},
			args: args{1, 99},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.IsBoundary(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("IsBoundary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deathBoundary_Copy(t *testing.T) {
	tests := []struct {
		name string
		b    deathBoundary
	}{
		{
			name: "should create a deep copy",
			b: deathBoundary{
				1: {1: wallExists{}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newB := tt.b.Copy()
			assert.Equal(t, tt.b, newB)

			tt.b.Add(99, 1293)
			assert.NotEqual(t, tt.b, newB)
		})
	}
}
