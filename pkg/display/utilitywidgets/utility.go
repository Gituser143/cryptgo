/*
Copyright © 2021 Bhargav SNV bhargavsnv100@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utilitywidgets

// Utility represents a utilty displayed in the UI
type Utility int

const (
	// None is when no utility is being displayed and the page is visible
	None Utility = iota
	// Help is when the help table is displayed
	Help
	// Portfolio is when the mini portfolio table is displayed
	Portfolio
	// Change is when the change interval table is displayed
	Change
	// Duration is when the graph duration interval table is displayed
	Duration
	// Currency is when the currency selection table is displayed
	Currency
)
