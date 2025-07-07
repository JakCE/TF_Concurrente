import { Routes } from '@angular/router';
import { MainLayoutComponent } from './components/main-layout/main-layout.component';

export const routes: Routes = [
    {
    path: '',
    component: MainLayoutComponent,
    children: [
      {
        path: '',
        loadComponent: () => import('./components/pages/form-prediccion/form-prediccion.component').then(m => m.FormPrediccionComponent)
      },
      {
        path: 'estadisticas',
        loadComponent: () => import('./components/pages/estadisticas/estadisticas.component').then(m => m.EstadisticasComponent)
      }
    ]
  },
  { path: '**', redirectTo: '' }
];
