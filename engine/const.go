package engine

// Text display description and logo.
const (
	headerTpl = `

                                                                            _/
     _/_/_/    _/_/    _/  _/_/    _/_/    _/_/_/      _/_/    _/  _/_/  _/_/_/_/    _/_/    _/  _/_/
  _/    _/  _/    _/  _/_/      _/_/_/_/  _/    _/  _/    _/  _/_/        _/      _/_/_/_/  _/_/
 _/    _/  _/    _/  _/        _/        _/    _/  _/    _/  _/          _/      _/        _/
  _/_/_/    _/_/    _/          _/_/_/  _/_/_/      _/_/    _/            _/_/    _/_/_/  _/
     _/                                _/
_/_/                                  _/

	Project: %s
	Score: %d
	Grade: %d
	Time: %s
	Issues: %d

	`
	metricsHeaderTpl = `>> %s Linter %s find:`
	summaryHeaderTpl = ` %s: %s`
	errorInfoTpl     = `  %s at line %d`
)
