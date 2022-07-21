import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { RouterModule } from '@angular/router';

import { NavigationComponent } from './navigation/navigation.component';
import { PanelComponent } from './panel/panel.component';

@NgModule({
  imports: [CommonModule, FontAwesomeModule, RouterModule],
  declarations: [NavigationComponent, PanelComponent],
  exports: [NavigationComponent, PanelComponent],
  providers: [],
})
export class ComponentsModule {}
